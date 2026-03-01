package io.clawdroid.core.domain.usecase

import io.clawdroid.core.domain.repository.ChatRepository
import io.mockk.mockk
import io.mockk.verify
import org.junit.jupiter.api.Test

class DisconnectChatUseCaseTest {

    private val repository = mockk<ChatRepository>(relaxed = true)
    private val useCase = DisconnectChatUseCase(repository)

    @Test
    fun `invoke delegates to repository disconnect`() {
        useCase()

        verify(exactly = 1) { repository.disconnect() }
    }
}
