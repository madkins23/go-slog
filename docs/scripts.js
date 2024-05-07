function hdlrChoice(form) {
    window.location.href = "/go-slog/handler/" + form.tag.value + ".html";
    return true;
}
function testChoice(form) {
    window.location.href = "/go-slog/test/" + form.tag.value + ".html";
    return true;
}
function otherChoice(form) {
    window.location.href = "/go-slog/" + form.tag.value;
    return true;
}
function toggleTableRow(id) {
    let elem = document.getElementById(id);
    elem.style.display = elem.style.display === "none" ? "table-row" : "none";
}
