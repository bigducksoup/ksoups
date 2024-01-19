import 'animate.css'
import {ref} from "vue";


export const useAnimation = (refs) => {

    const animationClass = ref('animate__fadeInUp')

    const speedCLass = ref('animate__slower')

    const SetAnimationClass = (className) => {
        forceRemove()
        animationClass.value = className
    }

    const animate = () => {
        refs.forEach((element,index)=>{
            setTimeout(()=>{
                element.value.classList.add('animate__animated')
                element.value.classList.add(animationClass.value)
                element.value.classList.add(speedCLass.value)
            },index*100)
        })
    }

    const remove = () => {
        refs.forEach((element,index)=>{
            setTimeout(()=>{
                element.value.classList.remove('animate__animated')
                element.value.classList.remove(animationClass.value)
                element.value.classList.remove(speedCLass.value)
            },index*100)
        })
    }

    const forceRemove = () => {
        refs.forEach((element,index)=>{
            element.value.classList.remove('animate__animated')
            element.value.classList.remove(animationClass.value)
            element.value.classList.remove(speedCLass.value)
        })
    }

    const forceAdd = () => {
        refs.forEach((element,index)=>{
            element.value.classList.add('animate__animated')
            element.value.classList.add(animationClass.value)
            element.value.classList.add(speedCLass.value)
        })
    }

    const DoAnimate = (animation,speed) => {
        forceRemove()
        animationClass.value = animation
        speedCLass.value = speed
        animate()
    }

    const SameTimeAnimate = (animation,speed) => {
        forceRemove()
        animationClass.value = animation
        speedCLass.value = speed
        requestAnimationFrame(()=>{
            forceAdd()
        })
    }

    const AddRefs = (ref) => {
        refs.push(ref)
    }

    const RemoveRefs = (ref) => {
        refs.splice(refs.indexOf(ref),1)
    }

    const RemoveAllRefs = () => {
        refs = []
    }


    return {DoAnimate,SetAnimationClass,AddRefs,RemoveRefs,RemoveAllRefs,SameTimeAnimate}


}


