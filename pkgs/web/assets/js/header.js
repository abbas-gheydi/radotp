fetch('/header/')
.then(res => res.text())
.then(text => {
    let oldelem = document.querySelector("script#add_navbar");
    let newelem = document.createElement("div");
    newelem.innerHTML = text;
    oldelem.parentNode.replaceChild(newelem,oldelem);
})