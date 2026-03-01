package io.clawdroid.core.data.local.converter

import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertNull
import org.junit.jupiter.api.Test

class ConvertersTest {

    private val converters = Converters()

    @Test
    fun `fromStringList null returns null`() {
        assertNull(converters.fromStringList(null))
    }

    @Test
    fun `toStringList null returns null`() {
        assertNull(converters.toStringList(null))
    }

    @Test
    fun `round-trip with non-empty list`() {
        val original = listOf("apple", "banana", "cherry")

        val json = converters.fromStringList(original)
        val restored = converters.toStringList(json)

        assertEquals(original, restored)
    }

    @Test
    fun `round-trip with empty list`() {
        val original = emptyList<String>()

        val json = converters.fromStringList(original)
        val restored = converters.toStringList(json)

        assertEquals(original, restored)
    }

    @Test
    fun `round-trip with single element`() {
        val original = listOf("single")

        val json = converters.fromStringList(original)
        val restored = converters.toStringList(json)

        assertEquals(original, restored)
    }

    @Test
    fun `handles strings with special characters`() {
        val original = listOf("hello world", "foo\"bar", "a,b,c")

        val json = converters.fromStringList(original)
        val restored = converters.toStringList(json)

        assertEquals(original, restored)
    }
}
