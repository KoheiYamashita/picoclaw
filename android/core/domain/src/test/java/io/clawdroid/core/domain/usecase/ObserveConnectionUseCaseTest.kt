package io.clawdroid.core.domain.usecase

import io.clawdroid.core.domain.model.ConnectionState
import io.clawdroid.core.domain.repository.ChatRepository
import io.mockk.every
import io.mockk.mockk
import kotlinx.coroutines.flow.MutableStateFlow
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test

class ObserveConnectionUseCaseTest {

    private val connectionFlow = MutableStateFlow(ConnectionState.DISCONNECTED)
    private val repository = mockk<ChatRepository> {
        every { connectionState } returns connectionFlow
    }
    private val useCase = ObserveConnectionUseCase(repository)

    @Test
    fun `invoke returns repository connectionState StateFlow`() {
        val result = useCase()

        assertEquals(connectionFlow, result)
    }

    @Test
    fun `returned flow reflects connection state changes`() {
        connectionFlow.value = ConnectionState.CONNECTED

        assertEquals(ConnectionState.CONNECTED, useCase().value)
    }
}
