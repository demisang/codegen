<template>
  <div class="path-items">
    <span v-for="(item, k) of pathItems" :key="item.path" :class="['dir-item', {'active': k !== pathItems.length - 1}]">
      <a v-if="k !== pathItems.length - 1" @click.prevent="selectDir(item.path)" href="#">
        {{ item.name }}
      </a>
      <span v-else>
        {{ item.name }}
      </span>
      /
    </span>

    <template v-if="currentDirItems.length">
      [sub dirs: <span v-for="item of currentDirItems" :key="item.path" class="dir-item">
      <a @click.prevent="selectDir(item.path)" href="#">
        {{ item.name }}
      </a>
      <span>&nbsp;</span>
    </span>]
    </template>
    <span v-else>[no sub dirs]</span>
  </div>
</template>

<script lang="ts">
import {defineComponent} from "vue"
import {getDirectories} from "@/api";

interface PathItem {
  path: string
  name: string
}

export default defineComponent({
  name: "SelectDir",
  data() {
    return {
      // pathItems: [] as PathItem[],
      currentDirItems: [] as PathItem[],
    }
  },
  props: {
    selected: {type: String, required: true},
  },
  computed: {
    pathItems():PathItem[] {
      let items = [{path: "", "name": "root"}]
      let fullPath: string[] = []
      this.selected.split("/").filter(value => value !== "").forEach((value: string) => {
        fullPath.push(value)
        items.push({
          path: fullPath.join("/"),
          name: value,
        })
      })

      this.fetchCurrentDirItems(this.selected)

      return items
    }
  },
  emits: {
    changed: (path: string) => {
      return true
    }
  },
  methods: {
    selectDir(path: string) {
      this.$emit("changed", path)
      this.fetchCurrentDirItems(path)
    },
    fetchCurrentDirItems(path: string) {
      getDirectories(path, (data, statusCode) => {
        this.currentDirItems = []
        data.forEach((value: string) => {
          this.currentDirItems.push({
            path: this.selected != "" ? this.selected + "/" + value : value,
            name: value,
          })
        })
      })
    },
  }
})
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.dir-item {
  //padding-right: 10px;
}
</style>
