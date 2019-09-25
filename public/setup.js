hljs.configure({ tabReplace: '  ' });

const go = new Go();
var gl

window.onload = function () {
    // linkRange('count','count-value')
    linkRange('gravity', 'gravity-value')
    // linkRange('size','size-value')

    //Go code background
    // fetch('main.go').then( res=> res.text()).then(res => {
    // 	let codeEl = document.getElementById("code")
    // 	codeEl.innerHTML = res.replace(/&/g,"&amp;")
    // 												.replace(/>/g,"&gt;")
    // 												.replace(/</g,"&lt;") 
    // 	hljs.highlightBlock(codeEl)
    // })
    

    WebAssembly.instantiateStreaming(fetch('main.wasm'), go.importObject)
    .then(res => {
        go.run(res.instance)
    })
}

function linkRange(id, idValue) {
    let El = document.getElementById(id)
    let valEl = document.getElementById(idValue)
    El.addEventListener("input", function () { valEl.innerHTML = El.value })
    valEl.innerHTML = El.value
}