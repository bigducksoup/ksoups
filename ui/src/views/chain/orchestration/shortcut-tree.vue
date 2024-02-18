<script setup>
import { NTree } from 'naive-ui'
import { computed,defineEmits } from 'vue';
import {useShortcutGroup} from '../../../hooks/shortcut.js'

const emits = defineEmits(['startDrag'])

const { shortcutGroup } = useShortcutGroup()

const tree = computed(() => {

    const options = []

    for(const field in shortcutGroup.value){

        const shortcuts = shortcutGroup.value[field]
        const children = []

        for(const shortcut of shortcuts){
            children.push({
                key: shortcut['id'],
                label: shortcut['name'],
                proto: shortcut
            })
        }

        options.push({
            key: field,
            label: field,
            children: children,
            disabled: true
        })
    }

    return options
})


const nodeProps = ({option}) => {
    return {
        onmousedown:(e) => {
            if (option.disabled)return
            emits('startDrag',e,option)
        }
    }
}


</script>

<template>
    <div class="pl-1 pr-1">
        <n-tree
        :data="tree"
        :node-props="nodeProps"
        expand-on-click
        >

        </n-tree>
    </div>
</template>