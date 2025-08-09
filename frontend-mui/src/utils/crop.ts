// Folosit pentru a decupa imaginea după react-easy-crop
export async function getCroppedImg(src: string, area: any): Promise<Blob> {
    const img = document.createElement("img");
    img.src = src;

    await new Promise(r => (img.onload = r));

    const canvas = document.createElement("canvas");
    const ctx = canvas.getContext("2d")!;
    canvas.width = area.width;
    canvas.height = area.height;

    ctx.drawImage(
        img,
        area.x,
        area.y,
        area.width,
        area.height,
        0,
        0,
        area.width,
        area.height
    );

    return new Promise(resolve => {
        canvas.toBlob(b => resolve(b!), "image/jpeg");
    });
}
