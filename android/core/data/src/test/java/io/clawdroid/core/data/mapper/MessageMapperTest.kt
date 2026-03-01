package io.clawdroid.core.data.mapper

import io.clawdroid.core.data.local.entity.MessageEntity
import io.clawdroid.core.data.remote.dto.WsOutgoing
import io.clawdroid.core.domain.model.ImageData
import io.clawdroid.core.domain.model.MessageSender
import io.clawdroid.core.domain.model.MessageStatus
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertNotNull
import org.junit.jupiter.api.Assertions.assertNull
import org.junit.jupiter.api.Assertions.assertTrue
import org.junit.jupiter.api.Nested
import org.junit.jupiter.api.Test

class MessageMapperTest {

    @Nested
    inner class ToDomain {

        @Test
        fun `converts entity to ChatMessage`() {
            val entity = MessageEntity(
                id = "abc",
                content = "Hello",
                sender = "USER",
                imagePathList = null,
                timestamp = 1000L,
                status = "SENT",
                messageType = null,
            )

            val result = MessageMapper.toDomain(entity)

            assertEquals("abc", result.id)
            assertEquals("Hello", result.content)
            assertEquals(MessageSender.USER, result.sender)
            assertEquals(1000L, result.timestamp)
            assertEquals(MessageStatus.SENT, result.status)
            assertTrue(result.images.isEmpty())
            assertNull(result.messageType)
        }

        @Test
        fun `parses imagePathList JSON`() {
            val json = """[{"path":"/data/img.jpg","width":100,"height":200}]"""
            val entity = MessageEntity(
                id = "1",
                content = "photo",
                sender = "USER",
                imagePathList = json,
                timestamp = 2000L,
                status = "SENT",
            )

            val result = MessageMapper.toDomain(entity)

            assertEquals(1, result.images.size)
            assertEquals("/data/img.jpg", result.images[0].path)
            assertEquals(100, result.images[0].width)
            assertEquals(200, result.images[0].height)
        }

        @Test
        fun `invalid JSON in imagePathList returns empty images`() {
            val entity = MessageEntity(
                id = "1",
                content = "test",
                sender = "USER",
                imagePathList = "not valid json",
                timestamp = 3000L,
                status = "SENT",
            )

            val result = MessageMapper.toDomain(entity)

            assertTrue(result.images.isEmpty())
        }

        @Test
        fun `maps AGENT sender`() {
            val entity = MessageEntity(
                id = "1",
                content = "reply",
                sender = "AGENT",
                imagePathList = null,
                timestamp = 4000L,
                status = "RECEIVED",
                messageType = "status",
            )

            val result = MessageMapper.toDomain(entity)

            assertEquals(MessageSender.AGENT, result.sender)
            assertEquals(MessageStatus.RECEIVED, result.status)
            assertEquals("status", result.messageType)
        }
    }

    @Nested
    inner class ToEntityFromWsOutgoing {

        @Test
        fun `converts WsOutgoing to entity with AGENT sender`() {
            val dto = WsOutgoing(content = "Hi there", type = null)

            val result = MessageMapper.toEntity(dto)

            assertNotNull(result.id)
            assertEquals("Hi there", result.content)
            assertEquals("AGENT", result.sender)
            assertNull(result.imagePathList)
            assertEquals("RECEIVED", result.status)
            assertNull(result.messageType)
        }

        @Test
        fun `preserves message type from dto`() {
            val dto = WsOutgoing(content = "status text", type = "status")

            val result = MessageMapper.toEntity(dto)

            assertEquals("status", result.messageType)
        }

        @Test
        fun `generates unique UUIDs`() {
            val dto = WsOutgoing(content = "test")
            val result1 = MessageMapper.toEntity(dto)
            val result2 = MessageMapper.toEntity(dto)

            assertTrue(result1.id != result2.id)
        }

        @Test
        fun `timestamp is set to current time`() {
            val before = System.currentTimeMillis()
            val result = MessageMapper.toEntity(WsOutgoing(content = "test"))
            val after = System.currentTimeMillis()

            assertTrue(result.timestamp in before..after)
        }
    }

    @Nested
    inner class ToEntityFromTextAndImages {

        @Test
        fun `creates entity with USER sender`() {
            val result = MessageMapper.toEntity("Hello", emptyList(), MessageStatus.SENDING)

            assertEquals("Hello", result.content)
            assertEquals("USER", result.sender)
            assertEquals("SENDING", result.status)
            assertNull(result.imagePathList)
        }

        @Test
        fun `serializes images as JSON`() {
            val images = listOf(ImageData("/data/photo.jpg", 640, 480))

            val result = MessageMapper.toEntity("pic", images, MessageStatus.SENDING)

            assertNotNull(result.imagePathList)
            assertTrue(result.imagePathList!!.contains("/data/photo.jpg"))
            assertTrue(result.imagePathList!!.contains("640"))
            assertTrue(result.imagePathList!!.contains("480"))
        }

        @Test
        fun `empty images produces null imagePathList`() {
            val result = MessageMapper.toEntity("text", emptyList(), MessageStatus.SENT)

            assertNull(result.imagePathList)
        }

        @Test
        fun `generates unique UUIDs`() {
            val r1 = MessageMapper.toEntity("a", emptyList(), MessageStatus.SENDING)
            val r2 = MessageMapper.toEntity("b", emptyList(), MessageStatus.SENDING)

            assertTrue(r1.id != r2.id)
        }
    }

    @Nested
    inner class ToWsIncoming {

        @Test
        fun `creates WsIncoming with content`() {
            val result = MessageMapper.toWsIncoming("Hello", emptyList())

            assertEquals("Hello", result.content)
            assertNull(result.images)
            assertNull(result.inputMode)
        }

        @Test
        fun `passes images when non-empty`() {
            val images = listOf("base64data1", "base64data2")

            val result = MessageMapper.toWsIncoming("text", images)

            assertEquals(2, result.images!!.size)
            assertEquals("base64data1", result.images!![0])
        }

        @Test
        fun `empty images produces null`() {
            val result = MessageMapper.toWsIncoming("text", emptyList())

            assertNull(result.images)
        }

        @Test
        fun `passes inputMode`() {
            val result = MessageMapper.toWsIncoming("text", emptyList(), "voice")

            assertEquals("voice", result.inputMode)
        }
    }
}
