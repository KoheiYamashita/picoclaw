package io.clawdroid.core.data.repository

import io.clawdroid.core.data.local.ImageFileStorage
import io.clawdroid.core.data.local.dao.MessageDao
import io.clawdroid.core.data.local.entity.MessageEntity
import io.clawdroid.core.data.remote.WebSocketClient
import io.clawdroid.core.data.remote.dto.WsIncoming
import io.clawdroid.core.data.remote.dto.WsOutgoing
import io.clawdroid.core.domain.model.ConnectionState
import io.clawdroid.core.domain.model.ImageAttachment
import io.clawdroid.core.domain.model.ImageData
import io.mockk.coEvery
import io.mockk.coVerify
import io.mockk.every
import io.mockk.mockk
import io.mockk.slot
import io.mockk.verify
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.cancel
import kotlinx.coroutines.launch
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.flowOf
import kotlinx.coroutines.test.UnconfinedTestDispatcher
import kotlinx.coroutines.test.runTest
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertNull
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Nested
import org.junit.jupiter.api.Test

@OptIn(ExperimentalCoroutinesApi::class)
class ChatRepositoryImplTest {

    private val testDispatcher = UnconfinedTestDispatcher()
    private lateinit var repoScope: kotlinx.coroutines.CoroutineScope
    private lateinit var incomingMessages: MutableSharedFlow<WsOutgoing>
    private lateinit var connectionStateFlow: MutableStateFlow<ConnectionState>
    private lateinit var webSocketClient: WebSocketClient
    private lateinit var messageDao: MessageDao
    private lateinit var imageFileStorage: ImageFileStorage
    private lateinit var repository: ChatRepositoryImpl

    @BeforeEach
    fun setup() {
        repoScope = kotlinx.coroutines.CoroutineScope(testDispatcher)
        incomingMessages = MutableSharedFlow()
        connectionStateFlow = MutableStateFlow(ConnectionState.DISCONNECTED)

        webSocketClient = mockk<WebSocketClient>(relaxed = true)
        every { webSocketClient.incomingMessages } returns incomingMessages
        every { webSocketClient.connectionState } returns connectionStateFlow

        messageDao = mockk<MessageDao>(relaxed = true)
        every { messageDao.getRecentMessages(any()) } returns flowOf(emptyList())

        imageFileStorage = mockk<ImageFileStorage>()

        repository = ChatRepositoryImpl(webSocketClient, messageDao, repoScope, imageFileStorage)
    }

    @AfterEach
    fun tearDown() {
        repoScope.cancel()
    }

    @Test
    fun `connect delegates to webSocketClient`() {
        repository.connect()

        verify { webSocketClient.connect() }
    }

    @Test
    fun `disconnect delegates to webSocketClient`() {
        repository.disconnect()

        verify { webSocketClient.disconnect() }
    }

    @Test
    fun `loadMore increases display limit`() {
        // Trigger collection so the Lazily-started StateFlow subscribes to DAO
        repoScope.launch { repository.messages.collect {} }

        verify { messageDao.getRecentMessages(ChatRepositoryImpl.INITIAL_LOAD_COUNT) }

        repository.loadMore()
        verify { messageDao.getRecentMessages(ChatRepositoryImpl.INITIAL_LOAD_COUNT + ChatRepositoryImpl.PAGE_SIZE) }

        repository.loadMore()
        verify { messageDao.getRecentMessages(ChatRepositoryImpl.INITIAL_LOAD_COUNT + 2 * ChatRepositoryImpl.PAGE_SIZE) }
    }

    @Test
    fun `connectionState returns webSocketClient connectionState`() {
        assertEquals(connectionStateFlow, repository.connectionState)
    }

    @Nested
    inner class SendMessage {

        @Test
        fun `sendMessage inserts entity and sends via websocket`() = runTest {
            val imageData = ImageData("/path/img.jpg", 100, 200)
            val saveResult = ImageFileStorage.SaveResult(imageData, "base64data")
            coEvery { imageFileStorage.saveFromUri("content://photo") } returns saveResult
            coEvery { webSocketClient.send(any<WsIncoming>()) } returns true

            repository.sendMessage("Hello", listOf(ImageAttachment("content://photo")))

            coVerify { messageDao.insert(any()) }
            coVerify { webSocketClient.send(any<WsIncoming>()) }
            coVerify { messageDao.update(match { it.status == "SENT" }) }
        }

        @Test
        fun `sendMessage marks as FAILED when websocket send fails`() = runTest {
            coEvery { webSocketClient.send(any<WsIncoming>()) } returns false

            repository.sendMessage("fail")

            coVerify { messageDao.update(match { it.status == "FAILED" }) }
        }

        @Test
        fun `sendMessage with no images sends empty base64 list`() = runTest {
            coEvery { webSocketClient.send(any<WsIncoming>()) } returns true

            repository.sendMessage("text only")

            val wsSlot = slot<WsIncoming>()
            coVerify { webSocketClient.send(capture(wsSlot)) }
            assertNull(wsSlot.captured.images)
        }
    }

    @Nested
    inner class IncomingMessageHandling {

        @Test
        fun `status message updates statusLabel`() = runTest {
            incomingMessages.emit(WsOutgoing(content = "Thinking...", type = "status"))

            assertEquals("Thinking...", repository.statusLabel.value)
        }

        @Test
        fun `status_end clears statusLabel`() = runTest {
            incomingMessages.emit(WsOutgoing(content = "Thinking...", type = "status"))
            incomingMessages.emit(WsOutgoing(content = "", type = "status_end"))

            assertNull(repository.statusLabel.value)
        }

        @Test
        fun `normal message inserts entity and clears status`() = runTest {
            val entitySlot = slot<MessageEntity>()
            coEvery { messageDao.insert(capture(entitySlot)) } returns Unit

            incomingMessages.emit(WsOutgoing(content = "Hello!", type = null))

            assertEquals("Hello!", entitySlot.captured.content)
            assertEquals("AGENT", entitySlot.captured.sender)
            assertEquals("RECEIVED", entitySlot.captured.status)
        }

        @Test
        fun `exit type is ignored`() = runTest {
            incomingMessages.emit(WsOutgoing(content = "bye", type = "exit"))

            coVerify(exactly = 0) { messageDao.insert(any()) }
        }

        @Test
        fun `setup_required type is ignored`() = runTest {
            incomingMessages.emit(WsOutgoing(content = "", type = "setup_required"))

            coVerify(exactly = 0) { messageDao.insert(any()) }
        }

        @Test
        fun `tool_request invokes callback and sends response`() = runTest {
            val toolRequestJson = """{"request_id":"r1","action":"screenshot"}"""
            repository.onToolRequest = { "screenshot_result" }
            coEvery { webSocketClient.send(any<WsIncoming>()) } returns true

            incomingMessages.emit(WsOutgoing(content = toolRequestJson, type = "tool_request"))

            coVerify {
                webSocketClient.send(match<WsIncoming> {
                    it.type == "tool_response" && it.requestId == "r1" && it.content == "screenshot_result"
                })
            }
        }

        @Test
        fun `tool_request without callback sends fallback message`() = runTest {
            val toolRequestJson = """{"request_id":"r2","action":"tap"}"""
            repository.onToolRequest = null
            coEvery { webSocketClient.send(any<WsIncoming>()) } returns true

            incomingMessages.emit(WsOutgoing(content = toolRequestJson, type = "tool_request"))

            coVerify {
                webSocketClient.send(match<WsIncoming> {
                    it.content == "tool request handler not configured"
                })
            }
        }
    }
}
