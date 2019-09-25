function initWebGL() {
    const canvas = document.querySelector("#glCanvas");
    // Initialize the GL context
    const ctx = canvas.getContext("webgl");

    // Only continue if WebGL is available and working
    if (ctx === null) {
        alert("Unable to initialize WebGL. Your browser or machine may not support it.");
        return;
    }

    // Set clear color to black, fully opaque
    ctx.clearColor(0.0, 0.0, 0.0, 1.0);
    // Clear the color buffer with specified clear color
    ctx.clear(ctx.COLOR_BUFFER_BIT);

    const shaderProgram = initShaderProgram(ctx, vsSource, fsSource);

    const programInfo = {
        program: shaderProgram,
        attribLocations: {
            vertexPosition: ctx.getAttribLocation(shaderProgram, 'aVertexPosition'),
        },
        uniformLocations: {
            projectionMatrix: ctx.getUniformLocation(shaderProgram, 'uProjectionMatrix'),
            modelViewMatrix: ctx.getUniformLocation(shaderProgram, 'uModelViewMatrix'),
        },
    };

    const buffers = initBuffers(ctx)

    drawScene(ctx, programInfo, buffers)
}


function initBuffers(gl) {

    // Create a buffer for the square's positions.

    const positionBuffer = gl.createBuffer();

    // Select the positionBuffer as the one to apply buffer
    // operations to from here out.

    gl.bindBuffer(gl.ARRAY_BUFFER, positionBuffer);

    // Now create an array of positions for the square.

    const positions = [
        -1.0, 1.0,
        1.0, 1.0,
        -1.0, -1.0,
        1.0, -1.0,
    ];

    // Now pass the list of positions into WebGL to build the
    // shape. We do this by creating a Float32Array from the
    // JavaScript array, then use it to fill the current buffer.

    gl.bufferData(gl.ARRAY_BUFFER,
        new Float32Array(positions),
        gl.STATIC_DRAW);

    return {
        position: positionBuffer,
    };
}