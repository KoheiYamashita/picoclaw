package io.clawdroid.core.domain.usecase

import io.clawdroid.core.domain.model.ImageAttachment
import io.clawdroid.core.domain.repository.ChatRepository
import io.mockk.coEvery
import io.mockk.coVerify
import io.mockk.mockk
import kotlinx.coroutines.test.runTest
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.assertThrows

class SendMessageUseCaseTest {

    private val repository = mockk<ChatRepository>(relaxed = true)
    private val useCase = SendMessageUseCase(repository)

    @Test
    fun `invoke delegates to repository sendMessage`() = runTest {
        useCase("Hello")

        coVerify { repository.sendMessage("Hello", emptyList(), null) }
    }

    @Test
    fun `invoke passes images to repository`() = runTest {
        val images = listOf(ImageAttachment("content://photo1", "image/png"))

        useCase("Hi", images)

        coVerify { repository.sendMessage("Hi", images, null) }
    }

    @Test
    fun `invoke passes inputMode to repository`() = runTest {
        useCase("test", inputMode = "voice")

        coVerify { repository.sendMessage("test", emptyList(), "voice") }
    }

    @Test
    fun `invoke propagates exception from repository`() = runTest {
        coEvery { repository.sendMessage(any(), any(), any()) } throws RuntimeException("Network error")

        assertThrows<RuntimeException> {
            useCase("fail")
        }
    }
}
