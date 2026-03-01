package io.clawdroid.backend.api

import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test

class BackendStateTest {

    @Test
    fun `all enum values exist`() {
        val values = BackendState.entries

        assertEquals(4, values.size)
    }

    @Test
    fun `valueOf returns correct values`() {
        assertEquals(BackendState.STOPPED, BackendState.valueOf("STOPPED"))
        assertEquals(BackendState.STARTING, BackendState.valueOf("STARTING"))
        assertEquals(BackendState.RUNNING, BackendState.valueOf("RUNNING"))
        assertEquals(BackendState.ERROR, BackendState.valueOf("ERROR"))
    }
}
