
export const  toHtml = (text) => {

    let one = text.replace(/\n/g, "<br>");

    return one.replace(/ /g, '&nbsp&nbsp&nbsp')
}