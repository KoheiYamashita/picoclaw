package io.clawdroid.feature.chat

import io.clawdroid.core.domain.model.ChatMessage
import io.clawdroid.core.domain.model.ConnectionState
import io.clawdroid.core.domain.model.ImageAttachment
import io.clawdroid.core.domain.model.MessageSender
import io.clawdroid.core.domain.model.MessageStatus
import io.clawdroid.core.domain.usecase.ConnectChatUseCase
import io.clawdroid.core.domain.usecase.DisconnectChatUseCase
import io.clawdroid.core.domain.usecase.LoadMoreMessagesUseCase
import io.clawdroid.core.domain.usecase.ObserveConnectionUseCase
import io.clawdroid.core.domain.usecase.ObserveMessagesUseCase
import io.clawdroid.core.domain.usecase.ObserveStatusUseCase
import io.clawdroid.core.domain.usecase.SendMessageUseCase
import io.clawdroid.feature.chat.voice.VoiceModeManager
import io.clawdroid.feature.chat.voice.VoiceModeState
import io.mockk.coEvery
import io.mockk.coVerify
import io.mockk.every
import io.mockk.mockk
import io.mockk.verify
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.test.UnconfinedTestDispatcher
import kotlinx.coroutines.test.resetMain
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.test.setMain
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertNull
import org.junit.jupiter.api.Assertions.assertTrue
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Nested
import org.junit.jupiter.api.Test

@OptIn(ExperimentalCoroutinesApi::class)
class ChatViewModelTest {

    private val testDispatcher = UnconfinedTestDispatcher()

    private val messagesFlow = MutableStateFlow<List<ChatMessage>>(emptyList())
    private val connectionFlow = MutableStateFlow(ConnectionState.DISCONNECTED)
    private val statusFlow = MutableStateFlow<String?>(null)
    private val voiceModeStateFlow = MutableStateFlow(VoiceModeState())

    private val sendMessage = mockk<SendMessageUseCase>(relaxed = true)
    private val observeMessages = mockk<ObserveMessagesUseCase> {
        every { this@mockk() } returns messagesFlow
    }
    private val observeConnection = mockk<ObserveConnectionUseCase> {
        every { this@mockk() } returns connectionFlow
    }
    private val observeStatus = mockk<ObserveStatusUseCase> {
        every { this@mockk() } returns statusFlow
    }
    private val loadMoreMessages = mockk<LoadMoreMessagesUseCase>(relaxed = true)
    private val connectChat = mockk<ConnectChatUseCase>(relaxed = true)
    private val disconnectChat = mockk<DisconnectChatUseCase>(relaxed = true)
    private val voiceModeManager = mockk<VoiceModeManager>(relaxed = true) {
        every { state } returns voiceModeStateFlow
    }

    private lateinit var viewModel: ChatViewModel

    @BeforeEach
    fun setup() {
        Dispatchers.setMain(testDispatcher)
        viewModel = ChatViewModel(
            sendMessage, observeMessages, observeConnection, observeStatus,
            loadMoreMessages, connectChat, disconnectChat, voiceModeManager,
        )
    }

    @AfterEach
    fun tearDown() {
        Dispatchers.resetMain()
    }

    @Test
    fun `init calls connectChat`() {
        verify { connectChat() }
    }

    @Test
    fun `init collects messages into uiState`() {
        val msg = ChatMessage("1", "Hi", MessageSender.USER, emptyList(), 1000L, MessageStatus.SENT)
        messagesFlow.value = listOf(msg)

        assertEquals(1, viewModel.uiState.value.messages.size)
        assertEquals("Hi", viewModel.uiState.value.messages[0].content)
    }

    @Test
    fun `init collects connectionState into uiState`() {
        connectionFlow.value = ConnectionState.CONNECTED

        assertEquals(ConnectionState.CONNECTED, viewModel.uiState.value.connectionState)
    }

    @Test
    fun `init collects statusLabel into uiState`() {
        statusFlow.value = "Thinking..."

        assertEquals("Thinking...", viewModel.uiState.value.statusLabel)
    }

    @Test
    fun `init collects voiceModeState into uiState`() {
        voiceModeStateFlow.value = VoiceModeState(isActive = true)

        assertTrue(viewModel.uiState.value.voiceModeState.isActive)
    }

    @Nested
    inner class OnEvent {

        @Test
        fun `OnInputChanged updates inputText`() {
            viewModel.onEvent(ChatEvent.OnInputChanged("Hello"))

            assertEquals("Hello", viewModel.uiState.value.inputText)
        }

        @Test
        fun `OnSendClick sends message and clears input`() = runTest {
            viewModel.onEvent(ChatEvent.OnInputChanged("Test message"))
            viewModel.onEvent(ChatEvent.OnSendClick)

            assertEquals("", viewModel.uiState.value.inputText)
            coVerify { sendMessage("Test message", emptyList()) }
        }

        @Test
        fun `OnSendClick does nothing when input is empty and no images`() = runTest {
            viewModel.onEvent(ChatEvent.OnInputChanged(""))
            viewModel.onEvent(ChatEvent.OnSendClick)

            coVerify(exactly = 0) { sendMessage(any(), any()) }
        }

        @Test
        fun `OnSendClick does nothing when input is only whitespace and no images`() = runTest {
            viewModel.onEvent(ChatEvent.OnInputChanged("   "))
            viewModel.onEvent(ChatEvent.OnSendClick)

            coVerify(exactly = 0) { sendMessage(any(), any()) }
        }

        @Test
        fun `OnSendClick sets error on exception`() = runTest {
            coEvery { sendMessage(any(), any()) } throws RuntimeException("Network error")
            viewModel.onEvent(ChatEvent.OnInputChanged("fail"))

            viewModel.onEvent(ChatEvent.OnSendClick)

            assertEquals("Network error", viewModel.uiState.value.error)
        }

        @Test
        fun `OnSendClick clears pendingImages`() = runTest {
            val image = ImageAttachment("content://photo")
            viewModel.onEvent(ChatEvent.OnImageAdded(image))
            viewModel.onEvent(ChatEvent.OnInputChanged("with image"))

            viewModel.onEvent(ChatEvent.OnSendClick)

            assertTrue(viewModel.uiState.value.pendingImages.isEmpty())
        }

        @Test
        fun `OnImageAdded appends to pendingImages`() {
            val image1 = ImageAttachment("content://photo1")
            val image2 = ImageAttachment("content://photo2")

            viewModel.onEvent(ChatEvent.OnImageAdded(image1))
            viewModel.onEvent(ChatEvent.OnImageAdded(image2))

            assertEquals(2, viewModel.uiState.value.pendingImages.size)
        }

        @Test
        fun `OnImageRemoved removes by index`() {
            viewModel.onEvent(ChatEvent.OnImageAdded(ImageAttachment("content://a")))
            viewModel.onEvent(ChatEvent.OnImageAdded(ImageAttachment("content://b")))
            viewModel.onEvent(ChatEvent.OnImageAdded(ImageAttachment("content://c")))

            viewModel.onEvent(ChatEvent.OnImageRemoved(1))

            assertEquals(2, viewModel.uiState.value.pendingImages.size)
            assertEquals("content://a", viewModel.uiState.value.pendingImages[0].uri)
            assertEquals("content://c", viewModel.uiState.value.pendingImages[1].uri)
        }

        @Test
        fun `OnLoadMore delegates to loadMoreMessages`() {
            viewModel.onEvent(ChatEvent.OnLoadMore)

            verify { loadMoreMessages() }
        }

        @Test
        fun `OnError sets error message`() {
            viewModel.onEvent(ChatEvent.OnError("Something went wrong"))

            assertEquals("Something went wrong", viewModel.uiState.value.error)
        }

        @Test
        fun `OnErrorDismissed clears error`() {
            viewModel.onEvent(ChatEvent.OnError("error"))
            viewModel.onEvent(ChatEvent.OnErrorDismissed)

            assertNull(viewModel.uiState.value.error)
        }

        @Test
        fun `OnVoiceModeStart delegates to voiceModeManager`() {
            viewModel.onEvent(ChatEvent.OnVoiceModeStart)

            verify { voiceModeManager.start(any()) }
        }

        @Test
        fun `OnVoiceModeStop delegates to voiceModeManager`() {
            viewModel.onEvent(ChatEvent.OnVoiceModeStop)

            verify { voiceModeManager.stop() }
        }

        @Test
        fun `OnVoiceModeInterrupt delegates to voiceModeManager`() {
            viewModel.onEvent(ChatEvent.OnVoiceModeInterrupt)

            verify { voiceModeManager.interrupt() }
        }

        @Test
        fun `OnVoiceCameraToggle delegates to voiceModeManager`() {
            viewModel.onEvent(ChatEvent.OnVoiceCameraToggle)

            verify { voiceModeManager.toggleCamera() }
        }
    }

}
