package io.clawdroid.receiver

import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.content.Context
import androidx.core.app.NotificationCompat
import io.clawdroid.backend.api.R

object NotificationHelper {

    private const val CHANNEL_ID = "clawdroid_messages"
    private const val NOTIFICATION_ID = 1001

    const val ASSISTANT_CHANNEL_ID = "clawdroid_assistant"

    const val GATEWAY_CHANNEL_ID = "clawdroid_gateway"

    fun createNotificationChannel(context: Context) {
        val manager = context.getSystemService(NotificationManager::class.java)

        val messageChannel = NotificationChannel(
            CHANNEL_ID,
            context.getString(io.clawdroid.R.string.notification_channel_messages),
            NotificationManager.IMPORTANCE_DEFAULT
        ).apply {
            description = context.getString(io.clawdroid.R.string.notification_channel_messages_desc)
        }
        manager.createNotificationChannel(messageChannel)

        val assistantChannel = NotificationChannel(
            ASSISTANT_CHANNEL_ID,
            context.getString(io.clawdroid.R.string.notification_channel_assistant),
            NotificationManager.IMPORTANCE_LOW
        ).apply {
            description = context.getString(io.clawdroid.R.string.notification_channel_assistant_desc)
            setShowBadge(false)
        }
        manager.createNotificationChannel(assistantChannel)

        val gatewayChannel = NotificationChannel(
            GATEWAY_CHANNEL_ID,
            context.getString(io.clawdroid.R.string.notification_channel_gateway),
            NotificationManager.IMPORTANCE_LOW
        ).apply {
            description = context.getString(io.clawdroid.R.string.notification_channel_gateway_desc)
            setShowBadge(false)
        }
        manager.createNotificationChannel(gatewayChannel)
    }

    fun showMessageNotification(context: Context, content: String) {
        val launchIntent = context.packageManager.getLaunchIntentForPackage(context.packageName)
        val pendingIntent = PendingIntent.getActivity(
            context, 0, launchIntent,
            PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE
        )

        val notification = NotificationCompat.Builder(context, CHANNEL_ID)
            .setSmallIcon(R.drawable.ic_notification)
            .setContentTitle(context.getString(io.clawdroid.R.string.notification_title))
            .setContentText(content.take(200))
            .setStyle(NotificationCompat.BigTextStyle().bigText(content.take(1000)))
            .setPriority(NotificationCompat.PRIORITY_DEFAULT)
            .setContentIntent(pendingIntent)
            .setAutoCancel(true)
            .build()

        val manager = context.getSystemService(NotificationManager::class.java)
        manager.notify(NOTIFICATION_ID, notification)
    }
}
