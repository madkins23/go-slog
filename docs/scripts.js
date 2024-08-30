function displayTable(id, display) {
    let elem = document.getElementById(id);
    elem.style.display = display;
}

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
function scoreChoice(form) {
    window.location.href = "/go-slog/scores/" + form.tag.value + "/summary.html";
    return true;
}
function toggleTableRow(id) {
    toggleElement(id, "table-row")
}
function toggleElement(id, on) {
    let elem = document.getElementById(id);
    elem.style.display = elem.style.display === "none" ? on : "none";
}
function checkboxElement(id, on) {
    let button = document.getElementById(id + '-checkbox')
    if (button != null) {
        const element = document.getElementById(id)
        if (element != null) {
            element.style.display = button.checked ? on : "none"
        }
    }
}
