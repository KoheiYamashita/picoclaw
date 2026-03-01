package io.clawdroid.core.domain.usecase

import io.clawdroid.core.domain.repository.ChatRepository
import io.mockk.mockk
import io.mockk.verify
import org.junit.jupiter.api.Test

class ConnectChatUseCaseTest {

    private val repository = mockk<ChatRepository>(relaxed = true)
    private val useCase = ConnectChatUseCase(repository)

    @Test
    fun `invoke delegates to repository connect`() {
        useCase()

        verify(exactly = 1) { repository.connect() }
    }
}
