import {createFile, getFileContent, modifyFile} from "../services/file.js";
import {useMessage} from "naive-ui";
import {createDir} from "../services/dir.js";
import {diffLines} from "../services/diff.js";

const useFSOperation = () => {

    const message = useMessage();

    /**
     * @param {Object} createFileForm
     * @param {String} createFileForm.probeId
     * @param {String} createFileForm.path string
     * @param {String} createFileForm.permission string
     * @return {Promise<Boolean>}
     * @constructor
     */
    const CreateFile = async (createFileForm) => {

        const probeId = createFileForm.probeId;
        const path = createFileForm.path;
        const perm = createFileForm.permission;

        const res = await createFile(probeId, path, perm);

        if (res.code !== 200){
            message.error(res.msg);
            return false;
        }
        message.success(res.msg);
        return true;
    }

    /**
     *
     * @param {Object} createDirForm
     * @param {String} createDirForm.probeId
     * @param {String} createDirForm.path
     * @param {String} createDirForm.permission
     * @return {Promise<boolean>}
     * @constructor
     */
    const CreateDir = async (createDirForm) => {

            const probeId = createDirForm.probeId;
            const path = createDirForm.path;
            const perm = createDirForm.permission;

            const res = await createDir(probeId, path, perm);

            if (res.code !== 200){
                message.error(res.msg);
                return false;
            }
            message.success(res.msg);
            return true;
    }


    const GetFileContent = async (probeId, path) => {

        let res = await getFileContent(probeId, path);

        if (res.code !== 200){
            message.error(res.msg);
            return  null
        }

        return res.data.content;
    }

    const EditFile = async (probeId, path, originContent,editedContent) => {

        console.log(originContent)
        console.log(editedContent)

        let diffRes = diffLines(originContent, editedContent);

        let res = await modifyFile(probeId, path, diffRes);

        if (res.code !== 200){
            message.error(res.msg);
            return false;
        }
        return true;
    }


    return {
        CreateFile,
        GetFileContent,
        CreateDir,
        EditFile
    }

}


export {
    useFSOperation
}
