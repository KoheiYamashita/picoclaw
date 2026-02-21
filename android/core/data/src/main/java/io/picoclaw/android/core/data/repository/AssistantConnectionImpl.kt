package io.picoclaw.android.core.data.repository

import io.ktor.client.HttpClient
import io.picoclaw.android.core.data.remote.WebSocketClient
import io.picoclaw.android.core.data.remote.dto.WsIncoming
import io.picoclaw.android.core.domain.model.AssistantMessage
import io.picoclaw.android.core.domain.model.ConnectionState
import io.picoclaw.android.core.domain.repository.AssistantConnection
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.cancel
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.SharedFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asSharedFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import java.util.UUID

class AssistantConnectionImpl(
    private val httpClient: HttpClient
) : AssistantConnection {

    private val scope = CoroutineScope(SupervisorJob() + Dispatchers.IO)
    private val clientId = UUID.randomUUID().toString()
    private val wsClient = WebSocketClient(httpClient, scope, clientId, "assistant")

    private val _messages = MutableSharedFlow<AssistantMessage>(extraBufferCapacity = 64)
    override val messages: SharedFlow<AssistantMessage> = _messages.asSharedFlow()

    private val _statusText = MutableStateFlow<String?>(null)
    override val statusText: StateFlow<String?> = _statusText.asStateFlow()

    override val connectionState: StateFlow<ConnectionState> = wsClient.connectionState

    init {
        scope.launch {
            wsClient.incomingMessages.collect { dto ->
                when (dto.type) {
                    "status" -> _statusText.value = dto.content
                    "status_end" -> _statusText.value = null
                    else -> {
                        _statusText.value = null
                        _messages.emit(AssistantMessage(content = dto.content, type = dto.type))
                    }
                }
            }
        }
    }

    override fun connect(wsUrl: String) {
        wsClient.wsUrl = wsUrl
        wsClient.connect()
    }

    override fun disconnect() {
        wsClient.disconnect()
        scope.cancel()
    }

    override suspend fun send(text: String, images: List<String>, inputMode: String) {
        val dto = WsIncoming(
            content = text,
            images = images.ifEmpty { null },
            inputMode = inputMode
        )
        wsClient.send(dto)
    }
}
