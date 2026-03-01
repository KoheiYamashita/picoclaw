plugins {
    alias(libs.plugins.android.application) apply false
    alias(libs.plugins.android.library) apply false
    alias(libs.plugins.compose.compiler) apply false
    alias(libs.plugins.serialization) apply false
    alias(libs.plugins.ksp) apply false
    alias(libs.plugins.ktlint) apply false
    alias(libs.plugins.kover) apply false
}

subprojects {
    apply(plugin = "org.jlleitschuh.gradle.ktlint")

    tasks.withType<Test> {
        useJUnitPlatform()
    }

    afterEvaluate {
        configurations.findByName("testRuntimeOnly")?.let {
            dependencies.add("testRuntimeOnly", libs.junit5.launcher.get().toString())
        }
    }
}
