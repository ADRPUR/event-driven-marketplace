// src/components/CropDialog.tsx
import { useState, useCallback } from "react";
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Slider,
    Button,
    Box,
} from "@mui/material";
import Cropper from "react-easy-crop";
import type { Area } from "react-easy-crop";

interface CropDialogProps {
    /** URL sau obiect-URL de la <input type="file"> */
    src: string | null;
    /** Ajunge aici când utilizatorul mută crop-ul (pentru previzualizare live) */
    onComplete: (_: Area, areaPixels: Area) => void;
    /** Click pe “Cancel” / închiderea dialogului */
    onCancel: () => void;
    /** Click pe “Save” – returnează zona crop-ată în coordonate px */
    onSave: () => void;
    /** Dialog deschis? (nu e obligatoriu, dar îl expun din motive de control) */
    open?: boolean;
}

export default function CropDialog({
                                       src,
                                       onComplete,
                                       onCancel,
                                       onSave,
                                       open = !!src,
                                   }: CropDialogProps) {
    const [zoom, setZoom] = useState(1);
    const [crop, setCrop] = useState<{ x: number; y: number }>({ x: 0, y: 0 });

    // react-easy-crop cere să-i furnizezi un onCropComplete => zona în px
    const handleCropComplete = useCallback(
        (_: Area, areaPixels: Area) => {
            onComplete(_, areaPixels);
        },
        [onComplete],
    );

    return (
        <Dialog maxWidth="sm" fullWidth open={open} onClose={onCancel}>
            <DialogTitle>Crop Photo</DialogTitle>
            <DialogContent dividers sx={{ height: 400, position: "relative", p: 0 }}>
                {src && (
                    <Cropper
                        image={src}
                        aspect={1}
                        cropShape="rect"
                        objectFit="horizontal-cover"
                        showGrid={false}
                        crop={crop}
                        onCropChange={setCrop}
                        zoom={zoom}
                        onZoomChange={setZoom}
                        onCropComplete={handleCropComplete}
                    />
                )}
            </DialogContent>

            {/* Zoom slider */}
            <Box px={3} mt={2}>
                <Slider
                    value={zoom}
                    min={-1}
                    max={3}
                    step={0.1}
                    onChange={(_, v) => setZoom(v as number)}
                />
            </Box>

            <DialogActions sx={{ px: 3, pb: 2 }}>
                <Button onClick={onCancel} variant="outlined">
                    Cancel
                </Button>
                <Button onClick={onSave} variant="contained">
                    Save
                </Button>
            </DialogActions>
        </Dialog>
    );
}

