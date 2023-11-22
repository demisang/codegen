<template>
  <div class="row">
    <div class="col-3">
      <div
          class="list-group"
          v-for="template in templates" v-bind:key="template.id"
      >
        <button
            type="button"
            @click="selectTemplate(template)"
            :class="[{'active': template.id === selectedTemplateId}, 'list-group-item', 'list-group-item-action']"
        >
          <h5 class="">{{ template.name }}</h5>
          <div class="hint">{{ template.description }}</div>
        </button>
      </div>
    </div>
    <div class="col-9">
      <form v-if="selectedTemplateId.length" method="post">
        <div class="mb-3">
          <label for="targetDir" class="form-label">Target directory</label>
          <SelectDir :selected="targetDir" @changed="changeTargetDir"></SelectDir>
          <!--          <input type="text" :value="targetDir" class="form-control" id="targetDir" aria-describedby="targetDirHelp">-->
          <div id="targetDirHelp" class="form-text">Select target directory</div>
        </div>

        <div
            v-for="(placeholder, key) in placeholders"
            v-bind:key="placeholder.value"
            :class="['row g-3 align-items-center placeholder-item', 'placeholder' + key]"
        >
          <div class="col-3 text-end">
            <label :for="placeholder.value" class="form-label">{{ placeholder.value }}</label>
          </div>
          <div class="col-auto">
            <input type="text" v-model="placeholder.replace" class="form-control" :id="placeholder.value"
                   :aria-describedby="placeholder.value + 'Help'">
          </div>
          <div class="col-auto">
            <div :id="placeholder.value + 'Help'" class="form-text">{{ placeholder.description }}</div>
          </div>
        </div>

        <div class="row mt-2">
          <div class="col">
            <button @click="showPreview()" type="button" class="btn btn-primary"><i class="bi bi-binoculars-fill"></i>
              Preview
            </button>
          </div>
          <div class="col text-end">
            <button @click="resetPlaceholders()" type="button" class="btn btn-secondary"><i class="bi bi-eraser"></i>
              Reset
            </button>
            &nbsp;
            <button @click="generate()" type="button" class="btn btn-success"><i class="bi bi-tools"></i> Generate
            </button>
          </div>
        </div>
      </form>
    </div>
  </div>

  <div class="row">
    <div class="col">
      <hr v-if="previewList.length || helpMessage.length"/>
    </div>
  </div>

  <div class="row">
    <div class="col-4">
      <div
          v-for="(item, key) of previewList" :key="item.path"
          @click="showFile(key)"
          :class="['preview-file', {'active': selectedFileIndex === key}]"
      >
        <template v-if="item.is_dir">
          <i :class="['bi', item.is_new ? 'bi-folder-plus' : 'bi-folder']"></i>
        </template>
        <template v-else>
          <i :class="['bi', item.is_new ? 'bi-file-plus' : 'bi-file']"></i>
          <!--          <i class="bi bi-folder-plus"></i>-->
        </template>
        root{{ targetDir.length ? "/" + targetDir : "" }}<span v-html="colorizedReplace(item.path)"></span>
      </div>
    </div>
    <div class="col-8">
      <div v-if="helpMessage.length" class="shadow p-3 mb-5 bg-body-tertiary rounded">
        <pre>{{ helpMessage }}</pre>
      </div>
      <pre class="file-content-preview" v-html="selectedFileContent"></pre>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent} from "vue"
import Template from "@/models/template"
import {generate, getRawList, getTemplates} from "@/api";
import Placeholder from "@/models/placeholder";
import SelectDir from "@/components/SelectDir.vue";
import PreviewItem from "@/models/preview_item";

const escapeHtml = (unsafe: string) => {
  return unsafe.replaceAll('&', '&amp;').replaceAll('<', '&lt;')
      .replaceAll('>', '&gt;').replaceAll('"', '&quot;')
      .replaceAll("'", '&#039;');
}

let originalTemplates: Template[] = []

export default defineComponent({
  name: "TemplatesList",
  components: {SelectDir},
  data() {
    return {
      selectedTemplateId: "",
      selectedFileIndex: -1,
      helpMessage: "",
      // selectedFileContent: "",
      placeholders: [] as Placeholder[],
      targetDir: "",
      templates: [] as Template[],
      // originalTemplates: [] as Template[],
      previewList: [] as PreviewItem[],
    }
  },
  computed: {
    selectedFileContent(): string {
      let result = ""

      if (!(this.selectedFileIndex in this.previewList) || !this.previewList[this.selectedFileIndex].content.length) {
        return result
      }

      result = escapeHtml(this.previewList[this.selectedFileIndex].content)
      result = "<i class='line-number'></i>" + result
      result = result.replaceAll(/\n/g, "\n<i class='line-number'></i>")
      result = this.colorizedReplace(result)

      return result
    },
  },
  mounted() {
    getTemplates("", (data, statusCode) => {
      this.templates = data
      originalTemplates = JSON.parse(JSON.stringify(data))
    })
  },
  methods: {
    getTemplateByID(templateId: string) {
      for (const template of this.templates) {
        if (template.id === templateId) {
          return template
        }
      }
    },
    selectTemplate(template: Template) {
      this.selectedTemplateId = template.id
      this.targetDir = template.target_dir
      this.placeholders = template.placeholders

      this.reset()
    },
    reset() {
      this.selectedFileIndex = -1
      this.previewList = []
      this.helpMessage = ""
    },
    resetPlaceholders() {
      const originalTemplate = originalTemplates.find((template) => template.id === this.selectedTemplateId)
      if (originalTemplate) {
        this.placeholders = JSON.parse(JSON.stringify(originalTemplate.placeholders))
      }
      this.showPreview()
    },
    changeTargetDir(path: string) {
      this.targetDir = path
      if (this.previewList.length) {
        this.showPreview()
      }
    },
    showPreview() {
      if (!this.selectedTemplateId.length) {
        return
      }

      this.helpMessage = ""
      getRawList(this.selectedTemplateId, this.targetDir, this.placeholders, (data: PreviewItem[]) => {
        this.previewList = data
      })
    },
    showFile(key: number) {
      this.selectedFileIndex = key
    },
    generate() {
      generate(this.selectedTemplateId, this.targetDir, this.placeholders, (helpMessage: string) => {
        this.helpMessage = helpMessage
      })
    },
    colorizedReplace(value: string): string {
      for (const key in this.placeholders) {
        const color = "color" + key
        let replace = escapeHtml(this.placeholders[key].replace)
        replace = `<span class='${color}'>${replace}</span>`
        value = value.replaceAll(this.placeholders[key].value, replace)
      }

      return value
    }
  }
})
</script>

<style lang="scss">
.placeholder-item {
  margin-bottom: 10px;
}

.preview-file {
  cursor: pointer;

  &:hover, &.active {
    background-color: rgb(43, 48, 53);
  }
}

.file-content-preview {
  tab-size: 4;
  counter-reset: lineNumber 0;

  .line-number {
    counter-increment: lineNumber 1;
    width: 28px;
    display: inline-block;
    text-align: right;
    margin-right: 10px;
    color: #555;

    &:before {
      content: counter(lineNumber);
    }
  }
}

$color1: #7fff00; // chartreuse
$color2: #00bfff; // deepskyblue
$color3: #ff4500; // orangered
$color4: #00fa9a; // mediumspringgreen
$color5: #ffa500; // orange
$color6: #ffff00; // yellow
$color7: #ff00ff; // fuchsia
$color8: #00ffff; // aqua
$color9: #4169e1; // royalblue
$color10: #ff1493; // deeppink
$color11: #0000ff; // blue
$color12: #ee82ee; // violet
$color13: #f0e68c; // khaki
$color14: #ffa07a; // lightsalmon
$color15: #2e8b57; // seagreen
$color16: #a52a2a; // brown
$color17: #000080; // navy
$color18: #800000; // maroon
$color19: #808000; // olive
$color20: #696969; // dimgray
$color21: #dcdcdc; // gainsboro

.color0, .placeholder0 input {
  color: $color1;
  border-color: $color1;
}

.color1, .placeholder1 input {
  color: $color2;
  border-color: $color2;
}

.color2, .placeholder2 input {
  color: $color3;
  border-color: $color3;
}

.color3, .placeholder3 input {
  color: $color4;
  border-color: $color4;
}

.color4, .placeholder4 input {
  color: $color5;
  border-color: $color5;
}

.color5, .placeholder5 input {
  color: $color6;
  border-color: $color6;
}

.color6, .placeholder6 input {
  color: $color7;
  border-color: $color7;
}

.color7, .placeholder7 input {
  color: $color8;
  border-color: $color8;
}

.color8, .placeholder8 input {
  color: $color9;
  border-color: $color9;
}

.color9, .placeholder9 input {
  color: $color10;
  border-color: $color10;
}

.color10, .placeholder10 input {
  color: $color11;
  border-color: $color11;
}

.color11, .placeholder11 input {
  color: $color12;
  border-color: $color12;
}

.color12, .placeholder12 input {
  color: $color13;
  border-color: $color13;
}

.color13, .placeholder13 input {
  color: $color14;
  border-color: $color14;
}

.color14, .placeholder14 input {
  color: $color15;
  border-color: $color15;
}

.color15, .placeholder15 input {
  color: $color16;
  border-color: $color16;
}

.color16, .placeholder16 input {
  color: $color17;
  border-color: $color17;
}

.color17, .placeholder17 input {
  color: $color18;
  border-color: $color18;
}

.color18, .placeholder18 input {
  color: $color19;
  border-color: $color19;
}

.color19, .placeholder19 input {
  color: $color20;
  border-color: $color20;
}

.color20, .placeholder20 input {
  color: $color21;
  border-color: $color21;
}
</style>
