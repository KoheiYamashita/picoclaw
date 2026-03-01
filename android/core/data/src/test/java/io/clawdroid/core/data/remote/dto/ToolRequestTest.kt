package io.clawdroid.core.data.remote.dto

import kotlinx.serialization.json.Json
import kotlinx.serialization.json.int
import kotlinx.serialization.json.jsonPrimitive
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertNotNull
import org.junit.jupiter.api.Assertions.assertNull
import org.junit.jupiter.api.Test

class ToolRequestTest {

    private val json = Json { ignoreUnknownKeys = true }

    @Test
    fun `deserialize with params`() {
        val jsonStr = """{"request_id":"req-1","action":"tap","params":{"x":100,"y":200}}"""

        val dto = json.decodeFromString<ToolRequest>(jsonStr)

        assertEquals("req-1", dto.requestId)
        assertEquals("tap", dto.action)
        assertNotNull(dto.params)
        assertEquals(100, dto.params!!["x"]!!.jsonPrimitive.int)
        assertEquals(200, dto.params!!["y"]!!.jsonPrimitive.int)
    }

    @Test
    fun `deserialize without params`() {
        val jsonStr = """{"request_id":"req-2","action":"screenshot"}"""

        val dto = json.decodeFromString<ToolRequest>(jsonStr)

        assertEquals("req-2", dto.requestId)
        assertEquals("screenshot", dto.action)
        assertNull(dto.params)
    }

    @Test
    fun `SerialName mapping for request_id`() {
        val jsonStr = """{"request_id":"r123","action":"test"}"""

        val dto = json.decodeFromString<ToolRequest>(jsonStr)

        assertEquals("r123", dto.requestId)
    }
}
