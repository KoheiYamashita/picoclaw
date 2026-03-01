package io.clawdroid.core.domain.usecase

import io.clawdroid.core.domain.repository.ChatRepository
import io.mockk.mockk
import io.mockk.verify
import org.junit.jupiter.api.Test

class LoadMoreMessagesUseCaseTest {

    private val repository = mockk<ChatRepository>(relaxed = true)
    private val useCase = LoadMoreMessagesUseCase(repository)

    @Test
    fun `invoke delegates to repository loadMore`() {
        useCase()

        verify(exactly = 1) { repository.loadMore() }
    }
}
