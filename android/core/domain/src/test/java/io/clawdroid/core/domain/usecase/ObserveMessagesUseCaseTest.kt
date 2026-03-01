package io.clawdroid.core.domain.usecase

import io.clawdroid.core.domain.model.ChatMessage
import io.clawdroid.core.domain.model.MessageSender
import io.clawdroid.core.domain.model.MessageStatus
import io.clawdroid.core.domain.repository.ChatRepository
import io.mockk.every
import io.mockk.mockk
import kotlinx.coroutines.flow.MutableStateFlow
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test

class ObserveMessagesUseCaseTest {

    private val messagesFlow = MutableStateFlow<List<ChatMessage>>(emptyList())
    private val repository = mockk<ChatRepository> {
        every { messages } returns messagesFlow
    }
    private val useCase = ObserveMessagesUseCase(repository)

    @Test
    fun `invoke returns repository messages StateFlow`() {
        val result = useCase()

        assertEquals(messagesFlow, result)
    }

    @Test
    fun `returned flow reflects repository state`() {
        val message = ChatMessage(
            id = "1",
            content = "Hello",
            sender = MessageSender.USER,
            timestamp = 1000L,
            status = MessageStatus.SENT,
        )
        messagesFlow.value = listOf(message)

        val result = useCase()

        assertEquals(1, result.value.size)
        assertEquals("Hello", result.value[0].content)
    }
}
