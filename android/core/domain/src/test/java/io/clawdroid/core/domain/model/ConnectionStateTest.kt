package io.clawdroid.core.domain.model

import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test

class ConnectionStateTest {

    @Test
    fun `all enum values exist`() {
        val values = ConnectionState.entries

        assertEquals(4, values.size)
    }

    @Test
    fun `valueOf returns correct values`() {
        assertEquals(ConnectionState.DISCONNECTED, ConnectionState.valueOf("DISCONNECTED"))
        assertEquals(ConnectionState.CONNECTING, ConnectionState.valueOf("CONNECTING"))
        assertEquals(ConnectionState.CONNECTED, ConnectionState.valueOf("CONNECTED"))
        assertEquals(ConnectionState.RECONNECTING, ConnectionState.valueOf("RECONNECTING"))
    }
}
