import * as Diff from 'diff'


const diffLines = (oldStr, newStr) => {
    return Diff.diffLines(oldStr, newStr)
}

export{
    diffLines
}