package io.clawdroid.backend.config

import io.mockk.coEvery
import io.mockk.coVerify
import io.mockk.mockk
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.UnconfinedTestDispatcher
import kotlinx.coroutines.test.advanceUntilIdle
import kotlinx.coroutines.test.resetMain
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.test.setMain
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonObject
import kotlinx.serialization.json.JsonPrimitive
import kotlinx.serialization.json.buildJsonObject
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertNull
import org.junit.jupiter.api.Assertions.assertTrue
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Nested
import org.junit.jupiter.api.Test

@OptIn(ExperimentalCoroutinesApi::class)
class ConfigViewModelTest {

    private val testDispatcher = UnconfinedTestDispatcher()
    private val apiClient = mockk<ConfigApiClient>(relaxed = true)
    private lateinit var viewModel: ConfigViewModel

    private val testSchema = ConfigSchema(
        sections = listOf(
            SchemaSection(
                key = "llm",
                label = "LLM Settings",
                fields = listOf(
                    SchemaField(key = "model", label = "Model", type = "string", secret = false),
                    SchemaField(key = "api_key", label = "API Key", type = "string", secret = true),
                ),
            ),
            SchemaSection(
                key = "gateway",
                label = "Gateway",
                fields = listOf(
                    SchemaField(key = "port", label = "Port", type = "int", secret = false),
                ),
            ),
        ),
    )

    private val testConfig: JsonObject = buildJsonObject {
        put("llm", buildJsonObject {
            put("model", JsonPrimitive("gpt-4"))
            put("api_key", JsonPrimitive("sk-123"))
        })
        put("gateway", buildJsonObject {
            put("port", JsonPrimitive(18790))
        })
    }

    @BeforeEach
    fun setup() {
        Dispatchers.setMain(testDispatcher)
    }

    @AfterEach
    fun tearDown() {
        Dispatchers.resetMain()
    }

    private fun createViewModel(): ConfigViewModel {
        return ConfigViewModel(apiClient)
    }

    @Nested
    inner class LoadData {

        @Test
        fun `loadData fetches schema and config successfully`() = runTest {
            coEvery { apiClient.getSchema() } returns testSchema
            coEvery { apiClient.getConfig() } returns testConfig

            viewModel = createViewModel()
            advanceUntilIdle()

            val listState = viewModel.uiState.value.listState
            assertTrue(listState is ListState.Loaded)
            val sections = (listState as ListState.Loaded).sections
            assertEquals(2, sections.size)
            assertEquals("llm", sections[0].key)
            assertEquals("LLM Settings", sections[0].label)
            assertEquals(2, sections[0].fieldCount)
        }

        @Test
        fun `loadData sets error state on failure`() = runTest {
            coEvery { apiClient.getSchema() } throws RuntimeException("Connection failed")

            viewModel = createViewModel()
            advanceUntilIdle()

            val listState = viewModel.uiState.value.listState
            assertTrue(listState is ListState.Error)
            assertEquals("Connection failed", (listState as ListState.Error).message)
        }

        @Test
        fun `loadData sets AuthRequired on AuthException`() = runTest {
            coEvery { apiClient.getSchema() } throws AuthException("HTTP 403: forbidden")

            viewModel = createViewModel()
            advanceUntilIdle()

            val listState = viewModel.uiState.value.listState
            assertTrue(listState is ListState.AuthRequired)
        }
    }

    @Nested
    inner class SectionSelection {

        @Test
        fun `onSectionSelected loads section fields`() = runTest {
            coEvery { apiClient.getSchema() } returns testSchema
            coEvery { apiClient.getConfig() } returns testConfig
            viewModel = createViewModel()
            advanceUntilIdle()

            viewModel.onSectionSelected("llm")

            val detail = viewModel.uiState.value.detailState
            assertEquals("llm", detail?.sectionKey)
            assertEquals("LLM Settings", detail?.sectionLabel)
            assertEquals(2, detail?.fields?.size)
            assertEquals("gpt-4", detail?.fields?.get(0)?.value)
        }

        @Test
        fun `onSectionSelected does not reload same section`() = runTest {
            coEvery { apiClient.getSchema() } returns testSchema
            coEvery { apiClient.getConfig() } returns testConfig
            viewModel = createViewModel()
            advanceUntilIdle()

            viewModel.onSectionSelected("llm")
            val state1 = viewModel.uiState.value.detailState
            viewModel.onSectionSelected("llm")
            val state2 = viewModel.uiState.value.detailState

            assertEquals(state1, state2)
        }
    }

    @Nested
    inner class FieldEditing {

        @Test
        fun `onFieldValueChanged updates field value`() = runTest {
            coEvery { apiClient.getSchema() } returns testSchema
            coEvery { apiClient.getConfig() } returns testConfig
            viewModel = createViewModel()
            advanceUntilIdle()
            viewModel.onSectionSelected("llm")

            viewModel.onFieldValueChanged("model", "claude-3")

            val field = viewModel.uiState.value.detailState?.fields?.find { it.key == "model" }
            assertEquals("claude-3", field?.value)
            assertEquals("gpt-4", field?.originalValue)
        }
    }

    @Nested
    inner class Save {

        @Test
        fun `onSave does nothing when no fields changed`() = runTest {
            coEvery { apiClient.getSchema() } returns testSchema
            coEvery { apiClient.getConfig() } returns testConfig
            viewModel = createViewModel()
            advanceUntilIdle()
            viewModel.onSectionSelected("llm")

            viewModel.onSave()

            coVerify(exactly = 0) { apiClient.saveConfig(any()) }
        }

        @Test
        fun `onSave sends changed fields`() = runTest {
            coEvery { apiClient.getSchema() } returns testSchema
            coEvery { apiClient.getConfig() } returns testConfig
            coEvery { apiClient.saveConfig(any()) } returns SaveConfigResult(status = "ok", restart = false)
            viewModel = createViewModel()
            advanceUntilIdle()
            viewModel.onSectionSelected("llm")
            viewModel.onFieldValueChanged("model", "claude-3")

            viewModel.onSave()
            advanceUntilIdle()

            coVerify { apiClient.saveConfig(any()) }
            val saveState = viewModel.uiState.value.saveState
            assertTrue(saveState is SaveState.Success)
        }

        @Test
        fun `onSave handles error response`() = runTest {
            coEvery { apiClient.getSchema() } returns testSchema
            coEvery { apiClient.getConfig() } returns testConfig
            coEvery { apiClient.saveConfig(any()) } returns SaveConfigResult(error = "validation failed")
            viewModel = createViewModel()
            advanceUntilIdle()
            viewModel.onSectionSelected("llm")
            viewModel.onFieldValueChanged("model", "")

            viewModel.onSave()
            advanceUntilIdle()

            val saveState = viewModel.uiState.value.saveState
            assertTrue(saveState is SaveState.Error)
            assertEquals("validation failed", (saveState as SaveState.Error).message)
        }

        @Test
        fun `onSave handles exception`() = runTest {
            coEvery { apiClient.getSchema() } returns testSchema
            coEvery { apiClient.getConfig() } returns testConfig
            coEvery { apiClient.saveConfig(any()) } throws RuntimeException("Network error")
            viewModel = createViewModel()
            advanceUntilIdle()
            viewModel.onSectionSelected("llm")
            viewModel.onFieldValueChanged("model", "new")

            viewModel.onSave()
            advanceUntilIdle()

            val saveState = viewModel.uiState.value.saveState
            assertTrue(saveState is SaveState.Error)
        }
    }

    @Nested
    inner class Navigation {

        @Test
        fun `onNavigateBackToList clears detailState and saveState`() = runTest {
            coEvery { apiClient.getSchema() } returns testSchema
            coEvery { apiClient.getConfig() } returns testConfig
            viewModel = createViewModel()
            advanceUntilIdle()
            viewModel.onSectionSelected("llm")

            viewModel.onNavigateBackToList()

            assertNull(viewModel.uiState.value.detailState)
            assertEquals(SaveState.Idle, viewModel.uiState.value.saveState)
        }

        @Test
        fun `dismissSaveResult resets saveState to Idle`() = runTest {
            coEvery { apiClient.getSchema() } returns testSchema
            coEvery { apiClient.getConfig() } returns testConfig
            coEvery { apiClient.saveConfig(any()) } returns SaveConfigResult(restart = true)
            viewModel = createViewModel()
            advanceUntilIdle()
            viewModel.onSectionSelected("llm")
            viewModel.onFieldValueChanged("model", "new")
            viewModel.onSave()
            advanceUntilIdle()

            viewModel.dismissSaveResult()

            assertEquals(SaveState.Idle, viewModel.uiState.value.saveState)
        }
    }
}
