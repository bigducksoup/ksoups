
<template>
    <div>
        <span @click="onClick('/',-1)" class="hover:underline hover:cursor-pointer" >@</span><span>/</span>
        <span v-for="(p,index) in pathes"> <span @click="onClick(p,index)" class="hover:underline hover:cursor-pointer" >{{ p }}</span><span>/</span></span>
    </div>
</template>


<script setup>
import { ref, watch } from 'vue';



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

    pathes.value = newPath.split('/').filter((p)=>p != '')

    console.log(pathes.value)

})

const onClick = (cur,index) =>{


    if(index == -1){
        emits('click','/')
        return
    }

    let clicked = '/'
    
    for(let i in pathes.value){
       clicked =  clicked.concat(pathes.value[i])

       if(i == index)break

       if(i != index){
        clicked =  clicked.concat('/')
       }
       
    }
    console.log(clicked)
    emits('click',clicked)

}

</script>