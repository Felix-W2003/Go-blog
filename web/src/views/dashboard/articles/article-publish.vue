<template>
  <div class="article-publish">
    <div class="title">
      <div class="left">
        <el-form :inline="true">
          <el-form-item label="文章标题">
            <el-input v-model="title" placeholder="请输入文章标题" clearable/>
          </el-form-item>
        </el-form>
      </div>
      <div class="right">
        <el-text>自动保存</el-text>
        <el-switch v-model="isAutoSaveEnabled"/>
        <el-button icon="Upload" @click="triggerImport">导入文章</el-button>
        <el-button type="danger" icon="Delete" @click="title='';text=''">清空文章</el-button>
        <el-button type="success" icon="Plus" @click="layoutStore.state.articleCreateVisible=true">发布文章</el-button>

        <!-- 导入 .md 用：真正的文件选择控件藏起来，由「导入文章」按钮触发 -->
        <input
            ref="fileInputRef"
            class="import-file-input"
            type="file"
            accept=".md,.markdown,text/markdown"
            @change="handleImportFile"
        />


        <!--  这里必须销毁，不然不会重新加载props-->
        <el-dialog
            v-model="articleCreateVisible"
            width="500"
            align-center
            destroy-on-close
            :before-close="articleCreateVisibleSynchronization"
        >
          <template #header>
            发布文章
          </template>
          <article-create-form
              :title="title"
              :content="text"
              :category="importedMeta.category"
              :tags="importedMeta.tags"
              :abstract="importedMeta.abstract"
              :cover="importedMeta.cover"
          />
          <template #footer>
          </template>
        </el-dialog>
      </div>
    </div>
    <MdEditor v-model="text" @onUploadImg="onUploadImg" @onSave="onSave" @onChange="onChange"/>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref, watch} from 'vue';
import {ElMessage, ElMessageBox} from 'element-plus';
import {MdEditor} from 'md-editor-v3';
import 'md-editor-v3/lib/style.css';
import axios from "axios";
import {useLayoutStore} from "@/stores/layout";
import ArticleCreateForm from "@/components/forms/ArticleCreateForm.vue";
import {parseMarkdownArticle} from "@/utils/markdown";
import type {AxiosResponse} from "axios";
import type {ApiResponse} from "@/utils/request";
import type {ImageUploadResponse} from "@/api/image";

const layoutStore = useLayoutStore()

const savedIsAutoSaveEnabled = localStorage.getItem('isAutoSaveEnabled');
const isAutoSaveEnabled = savedIsAutoSaveEnabled ? ref(savedIsAutoSaveEnabled === 'true') : ref(true)
watch(() => isAutoSaveEnabled.value, (newIsAutoSaveEnabled) => {
  localStorage.setItem('isAutoSaveEnabled', String(newIsAutoSaveEnabled))
})

const title = ref('')
const savedArticle = localStorage.getItem('article')
const text = savedArticle ? ref(savedArticle) : ref('')

/* =====================================================
   导入 Markdown 文章
   ===================================================== */

const fileInputRef = ref<HTMLInputElement>()

/**
 * 从 md 文件里解析出的元信息。
 * 标题与正文直接写进本页面（编辑器），这几项则通过 props 传给发布弹窗。
 */
const importedMeta = reactive({
  category: '',
  tags: [] as string[],
  abstract: '',
  cover: '',
})

const triggerImport = () => {
  fileInputRef.value?.click()
}

const handleImportFile = async (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]

  // 立刻清空，保证连续导入同一个文件也能再次触发 change
  input.value = ''

  if (!file) {
    return
  }

  if (!/\.(md|markdown)$/i.test(file.name)) {
    ElMessage.warning('目前只支持导入 .md / .markdown 文件')
    return
  }

  // 编辑器里已经有内容时先确认，避免手滑覆盖掉正在写的文章
  if (text.value.trim()) {
    try {
      await ElMessageBox.confirm(
          '导入会覆盖当前编辑器里的标题与正文，是否继续？',
          '提示',
          {
            confirmButtonText: '覆盖导入',
            cancelButtonText: '取消',
            type: 'warning',
          }
      )
    } catch {
      return
    }
  }

  const parsed = parseMarkdownArticle(await file.text(), file.name)

  title.value = parsed.title
  text.value = parsed.content

  importedMeta.category = parsed.category
  importedMeta.tags = parsed.tags
  importedMeta.abstract = parsed.abstract
  importedMeta.cover = parsed.cover

  // 尊重「自动保存」开关：导入后的正文也要能刷新找回
  if (isAutoSaveEnabled.value) {
    localStorage.setItem('article', text.value)
  }

  ElMessage.success(`已导入「${file.name}」`)
}

const onUploadImg = async (files: File[], callback: (urls: string[]) => void): Promise<void> => {
  const res = await Promise.all(
      files.map((file) => {
        return new Promise<AxiosResponse<ApiResponse<ImageUploadResponse>>>((resolve, reject) => {
          const form = new FormData();
          form.append('image', file);

          axios
              .post('/api/image/upload', form, {
                headers: {
                  'Content-Type': 'multipart/form-data',
                },
                withCredentials: true,
              })
              .then((response) => resolve(response))
              .catch((error) => reject(error));
        });
      })
  );

  callback(res.map((item) => item.data.data.url));
};

const onSave = (v: string, _: Promise<string>):void => {
  localStorage.setItem('article', v)
};

const onChange = (v: string):void => {
  if (isAutoSaveEnabled.value) {
    onSave(v,Promise.resolve(''))
  }
}

const articleCreateVisible = ref(layoutStore.state.articleCreateVisible);
watch(
    () => layoutStore.state.articleCreateVisible,
    (newValue) => {
      articleCreateVisible.value = newValue
    }
)

const articleCreateVisibleSynchronization = () => {
  layoutStore.state.articleCreateVisible = false
}
</script>

<style lang="scss">
.article-publish {
  height: 100%;

  // 导入文章用的隐藏文件选择框
  .import-file-input {
    display: none;
  }

  .title {
    display: flex;

    .left {
      .el-input {
        min-width: 400px;
      }
    }

    .right {
      margin-left: auto;

      .el-switch {
        margin-left: 20px;
        margin-right: 20px;
      }
    }
  }

  .md-editor {
    height: 95%;
  }
}
</style>

<style lang="scss">
.article-publish .md-editor .md-editor-toolbar-wrapper .md-editor-toolbar svg.md-editor-icon {
  height: 24px;
  width: 24px;
}
</style>
