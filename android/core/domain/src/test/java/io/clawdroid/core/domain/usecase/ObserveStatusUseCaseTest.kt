package io.clawdroid.core.domain.usecase

import io.clawdroid.core.domain.repository.ChatRepository
import io.mockk.every
import io.mockk.mockk
import kotlinx.coroutines.flow.MutableStateFlow
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertNull
import org.junit.jupiter.api.Test

class ObserveStatusUseCaseTest {

    private val statusFlow = MutableStateFlow<String?>(null)
    private val repository = mockk<ChatRepository> {
        every { statusLabel } returns statusFlow
    }
    private val useCase = ObserveStatusUseCase(repository)

    @Test
    fun `invoke returns repository statusLabel StateFlow`() {
        val result = useCase()

        assertEquals(statusFlow, result)
    }

    @Test
    fun `returned flow reflects null status initially`() {
        assertNull(useCase().value)
    }

    @Test
    fun `returned flow reflects status updates`() {
        statusFlow.value = "Thinking..."

        assertEquals("Thinking...", useCase().value)
    }
}
