package io.clawdroid.core.domain.model

import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertNotEquals
import org.junit.jupiter.api.Test

class ChatMessageTest {

    @Test
    fun `default images is empty list`() {
        val message = ChatMessage(
            id = "1",
            content = "Hello",
            sender = MessageSender.USER,
            timestamp = 1000L,
            status = MessageStatus.SENT,
        )

        assertEquals(emptyList<ImageData>(), message.images)
    }

    @Test
    fun `default messageType is null`() {
        val message = ChatMessage(
            id = "1",
            content = "test",
            sender = MessageSender.AGENT,
            timestamp = 2000L,
            status = MessageStatus.RECEIVED,
        )

        assertEquals(null, message.messageType)
    }

    @Test
    fun `equality based on all fields`() {
        val msg1 = ChatMessage("1", "Hello", MessageSender.USER, emptyList(), 1000L, MessageStatus.SENT)
        val msg2 = ChatMessage("1", "Hello", MessageSender.USER, emptyList(), 1000L, MessageStatus.SENT)

        assertEquals(msg1, msg2)
    }

    @Test
    fun `inequality when id differs`() {
        val msg1 = ChatMessage("1", "Hello", MessageSender.USER, emptyList(), 1000L, MessageStatus.SENT)
        val msg2 = ChatMessage("2", "Hello", MessageSender.USER, emptyList(), 1000L, MessageStatus.SENT)

        assertNotEquals(msg1, msg2)
    }

    @Test
    fun `copy preserves original values`() {
        val original = ChatMessage("1", "Hello", MessageSender.USER, emptyList(), 1000L, MessageStatus.SENT)
        val copied = original.copy(content = "Updated")

        assertEquals("Updated", copied.content)
        assertEquals("1", copied.id)
        assertEquals(MessageSender.USER, copied.sender)
    }

    @Test
    fun `MessageSender has USER and AGENT values`() {
        val values = MessageSender.entries

        assertEquals(2, values.size)
        assertEquals(MessageSender.USER, MessageSender.valueOf("USER"))
        assertEquals(MessageSender.AGENT, MessageSender.valueOf("AGENT"))
    }

    @Test
    fun `MessageStatus has all expected values`() {
        val values = MessageStatus.entries

        assertEquals(4, values.size)
        assertEquals(MessageStatus.SENDING, MessageStatus.valueOf("SENDING"))
        assertEquals(MessageStatus.SENT, MessageStatus.valueOf("SENT"))
        assertEquals(MessageStatus.FAILED, MessageStatus.valueOf("FAILED"))
        assertEquals(MessageStatus.RECEIVED, MessageStatus.valueOf("RECEIVED"))
    }
}
