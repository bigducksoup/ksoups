import {reactive, ref} from "vue";


export const useForm = (formValue, onSubmit) => {

    const initialForm = formValue

    let form = reactive(formValue)

    const Clear = () => {
        Object.keys(form).forEach(key => {
            form[key] = initialForm[key]
        })
    }

    const Submit = () => {
        onSubmit(form)
    }

    return {
        form,
        Submit,
        Clear
    }

}




export const useDataTable = (columns,data) => {

    let tableData = ref(data)

    let pagination = { pageSize: 10 }

    const AddItem = (item) => {
        tableData.value.push(item)
    }

    const RemoveItemByIndex = (index) => {
        tableData.value.splice(index, 1)
    }

    const SetData = (data) => {
        tableData.value = data
    }

    /**
     * @param {(value: unknown, index: number, array: unknown[]) => Boolean } condition
     */
    const RemoveItemByCondition =  (condition) => {
        tableData.value = tableData.value.filter(condition)
    }

    return { AddItem,RemoveItemByIndex, RemoveItemByCondition,SetData, tableData,columns}

}