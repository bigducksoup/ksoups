
<template>
    <div class="text-white h-full flex items-center">
        <div @click="onClick('/', -1)" class="hover:text-blue-500 hover:cursor-pointer h-full flex items-center">
            <n-icon size="19">
                <AtCircleOutline/>
            </n-icon>
        </div>

        <div v-for="(p, index) in pathes" class="inline-flex items-center">
            <div class="flex items-center mr-1 ml-1">
                <n-icon>
                    <ChevronForward />
                </n-icon>
            </div>
            <div @click="onClick(p, index)" class="hover:underline hover:text-blue-500 hover:cursor-pointer">
                {{ p }}
            </div>

        </div>
    </div>
</template>


<script setup>
import { ref, watch } from 'vue';
import { NIcon } from 'naive-ui';
import { ChevronForward , AtCircleOutline} from '@vicons/ionicons5';


const props = defineProps({
    'path': String
})

const emits = defineEmits(['click'])


const pathes = ref([])


watch(() => props.path, async (newPath, oldPath) => {



    if (newPath == '/') {
        pathes.value = []
        return
    }

    pathes.value = newPath.split('/').filter((p) => p != '')

    console.log(pathes.value)

})

const onClick = (cur, index) => {


    if (index == -1) {
        emits('click', '/')
        return
    }

    let clicked = '/'

    for (let i in pathes.value) {
        clicked = clicked.concat(pathes.value[i])

        if (i == index) break

        if (i != index) {
            clicked = clicked.concat('/')
        }

    }
    console.log(clicked)
    emits('click', clicked)

}

</script>