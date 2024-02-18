import { reactive, ref } from "vue";

export const useForm = (formValue, onSubmit) => {
  const initialForm = JSON.parse(JSON.stringify(formValue));

  let form = reactive(formValue);

  const Clear = () => {
    Object.keys(form).forEach((key) => {
      form[key] = initialForm[key];
    });
  };

  const Submit = () => {
    onSubmit(form);
  };

  return {
    form,
    Submit,
    Clear,
  };
};

export const useDataTable = (columns, data) => {
  let tableData = ref(data);

  let pagination = { pageSize: 10 };

  const AddItem = (item) => {
    tableData.value.push(item);
  };

  const RemoveItemByIndex = (index) => {
    tableData.value.splice(index, 1);
  };

  const SetData = (data) => {
    tableData.value = data;
  };

  /**
   * @param {(value: unknown, index: number, array: unknown[]) => Boolean } condition
   */
  const RemoveItemByCondition = (condition) => {
    tableData.value = tableData.value.filter(condition);
  };

  return {
    AddItem,
    RemoveItemByIndex,
    RemoveItemByCondition,
    SetData,
    tableData,
    columns,
  };
};

export const useMenu = (
  options = [],
  onUpdate = (key, item) => {
    console.log(key, item);
  }
) => {
  const MenuOption = ref(options);

  const Key = ref("");

  return { MenuOption, onUpdate, Key };
};

export const useSiderControl = () => {
  const collapsed = ref(false);

  /**
   * @type {Array<(mode:string)=>void>}
   */
  const onExpandHooks = [];

  /**
   * @type {Array<(mode:string)=>void>}
   */
  const onCollapseHooks = [];

  const expand = () => {
    collapsed.value = false;
    onExpandHooks.forEach((f) => {
      f.call("expand");
    });
  };

  const collapse = () => {
    collapsed.value = true;
    onCollapseHooks.forEach((f) => {
      f.call("collapse");
    });
  };

  /**
   * hookUpExpand
   * @param {(mode:string) => void} onExpand
   */
  const hookUpExpand = (onExpand) => {
    onExpandHooks.push(onExpand);
  };

  /**
   * hookUpCollapse
   * @param {(mode:string) => void} onExpand
   */
  const hookUpCollapse = (onCollapse) => {
    onCollapseHooks.push(onCollapse);
  };


  return {collapsed, expand, collapse, hookUpCollapse,hookUpExpand}


};
