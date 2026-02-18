package io.picoclaw.android.core.data.local

import android.content.Context
import android.graphics.BitmapFactory
import android.util.Base64
import io.picoclaw.android.core.domain.model.ImageData
import java.io.File
import java.util.UUID

class ImageFileStorage(context: Context) {

    private val imageDir = File(context.filesDir, "chat_images").also { it.mkdirs() }

    fun saveBase64ToFile(base64: String): ImageData {
        val bytes = Base64.decode(base64, Base64.DEFAULT)
        val file = File(imageDir, "${UUID.randomUUID()}.jpg")
        file.writeBytes(bytes)

        val opts = BitmapFactory.Options().apply { inJustDecodeBounds = true }
        BitmapFactory.decodeFile(file.absolutePath, opts)

        return ImageData(
            path = file.absolutePath,
            width = opts.outWidth,
            height = opts.outHeight
        )
    }
}
