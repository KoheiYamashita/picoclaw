package io.clawdroid.setup

import io.clawdroid.backend.api.GatewaySettingsStore
import io.ktor.client.HttpClient
import io.ktor.client.engine.okhttp.OkHttp
import io.ktor.client.request.header
import io.ktor.client.request.post
import io.ktor.client.request.put
import io.ktor.client.request.setBody
import io.ktor.client.statement.bodyAsText
import io.ktor.http.ContentType
import io.ktor.http.contentType
import io.ktor.http.isSuccess
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonObject
import kotlinx.serialization.json.jsonObject
import kotlinx.serialization.json.jsonPrimitive
import java.io.Closeable
import java.io.IOException

class SetupApiClient(private val settingsStore: GatewaySettingsStore) : Closeable {

    private val json = Json { ignoreUnknownKeys = true }
    private val client = HttpClient(OkHttp)

    private val baseUrl: String get() = settingsStore.settings.value.httpBaseUrl
    private val apiKey: String get() = settingsStore.settings.value.apiKey

    suspend fun init(body: JsonObject) {
        val response = client.post("$baseUrl/api/setup/init") {
            contentType(ContentType.Application.Json)
            setBody(body.toString())
        }
        if (!response.status.isSuccess()) {
            val errorMsg = parseError(response.bodyAsText())
            throw IOException("HTTP ${response.status.value}: $errorMsg")
        }
    }

    suspend fun complete(body: JsonObject) {
        val response = client.put("$baseUrl/api/setup/complete") {
            contentType(ContentType.Application.Json)
            setBody(body.toString())
            if (apiKey.isNotEmpty()) header("Authorization", "Bearer $apiKey")
        }
        if (!response.status.isSuccess()) {
            val errorMsg = parseError(response.bodyAsText())
            throw IOException("HTTP ${response.status.value}: $errorMsg")
        }
    }

    override fun close() {
        client.close()
    }

    private fun parseError(responseBody: String): String {
        return try {
            json.parseToJsonElement(responseBody).jsonObject["error"]?.jsonPrimitive?.content
                ?: "request failed"
        } catch (_: Exception) {
            "request failed"
        }
    }
}
