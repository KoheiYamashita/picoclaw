package io.clawdroid.core.data.remote.dto

import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertNull
import org.junit.jupiter.api.Assertions.assertTrue
import org.junit.jupiter.api.Test

class ToolResponseTest {

    private val json = Json { ignoreUnknownKeys = true }

    @Test
    fun `serialize success response`() {
        val dto = ToolResponse(requestId = "req-1", success = true, result = "done")

        val jsonStr = json.encodeToString(dto)
        val decoded = json.decodeFromString<ToolResponse>(jsonStr)

        assertEquals("req-1", decoded.requestId)
        assertTrue(decoded.success)
        assertEquals("done", decoded.result)
        assertNull(decoded.error)
    }

    @Test
    fun `serialize error response`() {
        val dto = ToolResponse(requestId = "req-2", success = false, error = "not found")

        val jsonStr = json.encodeToString(dto)
        val decoded = json.decodeFromString<ToolResponse>(jsonStr)

        assertEquals("req-2", decoded.requestId)
        assertEquals(false, decoded.success)
        assertNull(decoded.result)
        assertEquals("not found", decoded.error)
    }

    @Test
    fun `SerialName mapping for request_id`() {
        val jsonStr = """{"request_id":"r1","success":true}"""

        val dto = json.decodeFromString<ToolResponse>(jsonStr)

        assertEquals("r1", dto.requestId)
    }
}
