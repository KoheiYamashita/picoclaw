package io.picoclaw.android.assistant

import android.app.Notification
import android.content.Intent
import android.graphics.PixelFormat
import android.os.IBinder
import android.view.Gravity
import android.view.WindowManager
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.platform.ComposeView
import androidx.core.app.NotificationCompat
import androidx.lifecycle.LifecycleService
import androidx.lifecycle.setViewTreeLifecycleOwner
import androidx.savedstate.SavedStateRegistry
import androidx.savedstate.SavedStateRegistryController
import androidx.savedstate.SavedStateRegistryOwner
import androidx.savedstate.setViewTreeSavedStateRegistryOwner
import io.ktor.client.HttpClient
import io.picoclaw.android.core.data.remote.WebSocketClient
import io.picoclaw.android.core.data.repository.AssistantConnectionImpl
import io.picoclaw.android.core.domain.repository.AssistantConnection
import io.picoclaw.android.core.domain.repository.TtsSettingsRepository
import io.picoclaw.android.core.ui.theme.PicoClawTheme
import io.picoclaw.android.feature.chat.assistant.AssistantManager
import io.picoclaw.android.feature.chat.assistant.AssistantPillBar
import io.picoclaw.android.feature.chat.voice.CameraCaptureManager
import io.picoclaw.android.feature.chat.voice.SpeechRecognizerWrapper
import io.picoclaw.android.feature.chat.voice.TextToSpeechWrapper
import io.picoclaw.android.receiver.NotificationHelper
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.cancel
import org.koin.android.ext.android.inject

class AssistantService : LifecycleService(), SavedStateRegistryOwner {

    private val httpClient: HttpClient by inject()
    private val ttsSettingsRepo: TtsSettingsRepository by inject()

    private lateinit var serviceScope: CoroutineScope
    private lateinit var connection: AssistantConnection
    private lateinit var assistantManager: AssistantManager
    private lateinit var ttsWrapper: TextToSpeechWrapper
    private lateinit var sttWrapper: SpeechRecognizerWrapper
    private lateinit var cameraCaptureManager: CameraCaptureManager

    private var overlayView: ComposeView? = null
    private val windowManager by lazy { getSystemService(WINDOW_SERVICE) as WindowManager }

    private val savedStateRegistryController = SavedStateRegistryController.create(this)
    override val savedStateRegistry: SavedStateRegistry
        get() = savedStateRegistryController.savedStateRegistry

    override fun onCreate() {
        savedStateRegistryController.performAttach()
        savedStateRegistryController.performRestore(null)
        super.onCreate()

        serviceScope = CoroutineScope(SupervisorJob() + Dispatchers.Main)

        connection = AssistantConnectionImpl(httpClient)

        sttWrapper = SpeechRecognizerWrapper(this)
        ttsWrapper = TextToSpeechWrapper(this, ttsSettingsRepo.ttsConfig)
        cameraCaptureManager = CameraCaptureManager(this)

        assistantManager = AssistantManager(
            sttWrapper = sttWrapper,
            ttsWrapper = ttsWrapper,
            connection = connection,
            cameraCaptureManager = cameraCaptureManager,
            contentResolver = contentResolver
        )
    }

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        super.onStartCommand(intent, flags, startId)

        startForeground(NOTIFICATION_ID, buildNotification())

        // Resolve wsUrl from the main WebSocketClient
        val mainWsClient: WebSocketClient by inject()
        connection.connect(mainWsClient.wsUrl)

        addOverlay()
        assistantManager.start(serviceScope)

        return START_NOT_STICKY
    }

    override fun onDestroy() {
        removeOverlay()
        assistantManager.destroy()
        ttsWrapper.destroy()
        connection.disconnect()
        serviceScope.cancel()
        super.onDestroy()
    }

    override fun onBind(intent: Intent): IBinder? {
        super.onBind(intent)
        return null
    }

    private fun shutdown() {
        stopForeground(STOP_FOREGROUND_REMOVE)
        stopSelf()
    }

    private fun addOverlay() {
        if (overlayView != null) return

        val params = WindowManager.LayoutParams(
            WindowManager.LayoutParams.MATCH_PARENT,
            WindowManager.LayoutParams.WRAP_CONTENT,
            WindowManager.LayoutParams.TYPE_APPLICATION_OVERLAY,
            WindowManager.LayoutParams.FLAG_NOT_TOUCH_MODAL or
                WindowManager.LayoutParams.FLAG_LAYOUT_IN_SCREEN,
            PixelFormat.TRANSLUCENT
        ).apply {
            gravity = Gravity.BOTTOM
        }

        val view = ComposeView(this).apply {
            setViewTreeLifecycleOwner(this@AssistantService)
            setViewTreeSavedStateRegistryOwner(this@AssistantService)

            setContent {
                PicoClawTheme {
                    val state by assistantManager.state.collectAsState()
                    AssistantPillBar(
                        state = state,
                        onClose = { shutdown() },
                        onInterrupt = { assistantManager.interrupt() },
                        onCameraToggle = { assistantManager.toggleCamera() },
                        cameraCaptureManager = cameraCaptureManager
                    )
                }
            }
        }

        windowManager.addView(view, params)
        overlayView = view
    }

    private fun removeOverlay() {
        overlayView?.let {
            windowManager.removeView(it)
            overlayView = null
        }
    }

    private fun buildNotification(): Notification {
        return NotificationCompat.Builder(this, NotificationHelper.ASSISTANT_CHANNEL_ID)
            .setSmallIcon(android.R.drawable.ic_btn_speak_now)
            .setContentTitle("PicoClaw Assistant")
            .setContentText("Listening...")
            .setPriority(NotificationCompat.PRIORITY_LOW)
            .setOngoing(true)
            .build()
    }

    companion object {
        private const val NOTIFICATION_ID = 2001
    }
}
