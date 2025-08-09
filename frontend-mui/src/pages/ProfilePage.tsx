import {
    Box,
    Card,
    CardContent,
    Typography,
    Button,
    Divider,
    Stack,
    CircularProgress,
    Chip,
    TextField,
} from "@mui/material";
import {Save, Upload as UploadIcon} from "@mui/icons-material";
import {useCallback, useEffect, useState} from "react";
import type {Area} from "react-easy-crop";
import dayjs, {Dayjs} from "dayjs";

import {useAuthStore} from "../store/authStore";
import {getProfile, updateProfile, uploadPhoto} from "../api/user";
import {API_URL} from "../api/auth";
import type {UserDetails, Address} from "../types/user";

import CropDialog from "../components/CropDialog";
import {createImage} from "../utils/helper";   // helper ce încarcă imaginea în <img/>

/* ---------------- utilitare ---------------- */
const toISO = (d?: string | Dayjs) =>
    typeof d === "string" ? d.slice(0, 10) : d ? d.format("YYYY-MM-DD") : "";

/* decupează în canvas zona selectată */
async function getCroppedBlob(src: string, crop: Area): Promise<Blob> {
    const img = await createImage(src);
    const canvas = document.createElement("canvas");
    const ctx = canvas.getContext("2d")!;
    canvas.width = crop.width;
    canvas.height = crop.height;
    ctx.drawImage(
        img,
        crop.x,
        crop.y,
        crop.width,
        crop.height,
        0,
        0,
        crop.width,
        crop.height
    );
    return new Promise((resolve) =>
        canvas.toBlob((b) => b && resolve(b), "image/jpeg")
    );
}

/* ---------------- componentă principală ---------------- */
export default function ProfilePage() {
    const {token, login} = useAuthStore();

    const [data, setData] = useState<UserDetails | null>(null);
    const [edit, setEdit] = useState(false);
    const [saving, setSaving] = useState(false);
    const [uploading, setUploading] = useState(false);

    /* form „About” */
    const [form, setForm] = useState({
        firstName: "",
        lastName: "",
        phone: "",
        dateOfBirth: "",
        address: {
            line: "",
            city: "",
            postal_code: "",
            country: "",
        },
    });

    /* --- fetch profil --- */
    useEffect(() => {
        if (!token) return;
        (async () => {
            const p = await getProfile(token);
            setData(p);
            setForm({
                firstName: p.details.FirstName ?? "",
                lastName: p.details.LastName ?? "",
                phone: p.details.Phone ?? "",
                dateOfBirth: toISO(p.details.DateOfBirth),
                address: {
                    line: (p.details.Address as Address)?.line ?? "",
                    city: (p.details.Address as Address)?.city ?? "",
                    postal_code: (p.details.Address as Address)?.postal_code ?? "",
                    country: (p.details.Address as Address)?.country ?? "",
                }
            });
        })();
    }, [token]);

    /* --- upload + crop --- */
    const [cropOpen, setCropOpen] = useState(false);
    const [cropSrc, setCropSrc] = useState<string>();
    const [croppedArea, setCroppedArea] = useState<Area | null>(null);

    const handleSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (!e.target.files?.length) return;
        setCropSrc(URL.createObjectURL(e.target.files[0]));
        setCropOpen(true);
    };

    const handleCropComplete = useCallback(
        (_: Area, areaPixels: Area) => setCroppedArea(areaPixels),
        []
    );

    const handleCropSave = async () => {
        if (!cropSrc || !croppedArea || !token) return;
        setUploading(true);
        try {
            const blob = await getCroppedBlob(cropSrc, croppedArea);
            await uploadPhoto(token, new File([blob], "avatar.jpg"));
            const upd = await getProfile(token);
            setData(upd);
            login(upd, token);
            setCropOpen(false);
        } finally {
            setUploading(false);
        }
    };

    /* --- save About --- */
    const handleSave = async () => {
        if (!token) return;
        setSaving(true);
        try {
            await updateProfile(token, {
                firstName: form.firstName,
                lastName: form.lastName,
                phone: form.phone,
                dateOfBirth: form.dateOfBirth,
                address: form.address ? {
                    line: form.address.line,
                    city: form.address.city,
                    postal_code: form.address.postal_code,
                    country: form.address.country,
                } : null,
            });
            const upd = await getProfile(token);
            setData(upd);
            login(upd, token);
            setEdit(false);
        } finally {
            setSaving(false);
        }
    };

    if (!data)
        return (
            <Box p={5} textAlign="center">
                <CircularProgress/>
            </Box>
        );

    /* ---------------- UI ---------------- */
    return (
        <Box
            sx={{
                display: "grid",
                gap: 3,
                gridTemplateColumns: {xs: "1fr", md: "360px 1fr"}, // col. laterală + col. principală
            }}
        >
            {/* --------- card profil lateral --------- */}
            <Card sx={{width: "100%", pt: 2}}>
                <CardContent>
                    <Stack alignItems="center" spacing={1}>
                        {/* imagine decupată 4:5 */}
                        <Box
                            component="img"
                            src={
                                data.details.PhotoPath
                                    ? `${API_URL}${data.details.PhotoPath}?t=${Date.now()}`
                                    : undefined
                            }
                            alt={data.details.FirstName}
                            sx={{
                                width: 300,
                                height: 375,
                                objectFit: "cover",
                                borderRadius: 2,
                                boxShadow: 3,
                            }}
                        />
                        <Button
                            component="label"
                            variant="outlined"
                            size="small"
                            startIcon={<UploadIcon fontSize="small"/>}
                            disabled={uploading}
                            sx={{textTransform: "none", mt: 1}}
                        >
                            Change photo
                            <input hidden type="file" accept="image/*" onChange={handleSelect}/>
                        </Button>

                        <Typography variant="h6" mt={1}>
                            {data.details.FirstName} {data.details.LastName}
                        </Typography>
                        <Typography color="text.secondary" variant="body2">
                            {data.role}
                        </Typography>

                        <Divider sx={{width: "100%", my: 1}}/>

                        <Typography color="text.secondary" variant="caption">
                            Member since
                        </Typography>
                        <Typography variant="body2" fontWeight={600}>
                            {dayjs(data.details.CreatedAt).format("MMM DD, YYYY")}
                        </Typography>

                        <Divider sx={{width: "100%", my: 2}}/>
                        <Chip label="Active" color="success" size="small"/>
                    </Stack>
                </CardContent>
            </Card>

            {/* --------- coloană principală --------- */}
            <Box sx={{display: "grid", gap: 3}}>
                {/* ABOUT */}
                {!edit && (
                    <Box
                        sx={{
                            display: "grid",
                            gap: 3,
                            gridTemplateColumns: {xs: "1fr", md: "1fr 1fr"},
                        }}
                    >
                        <Card>
                            <CardContent>
                                <Typography variant="subtitle1" fontWeight={600} mb={2}>
                                    About
                                </Typography>

                                <Button
                                    size="small"
                                    variant="outlined"
                                    onClick={() => setEdit(true)}
                                    sx={{marginBottom: 2, textTransform: "none"}}
                                >
                                    Edit
                                </Button>

                                <InfoRow label="First Name" value={data.details.FirstName}/>
                                <InfoRow label="Last Name" value={data.details.LastName}/>
                                <InfoRow label="Contact No." value={data.details.Phone}/>
                                <InfoRow label="Email" value={data.email}/>
                                <InfoRow label="Birthday"
                                         value={dayjs(data.details.DateOfBirth).format("MMM DD, YYYY")}/>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardContent>
                                <Typography variant="subtitle1" fontWeight={600} mb={2}>
                                    Address
                                </Typography>
                                <InfoRow label="Street Line" value={(data.details.Address as Address)?.line}/>
                                <InfoRow label="City" value={(data.details.Address as Address)?.city}/>
                                <InfoRow label="Postal Code" value={(data.details.Address as Address)?.postal_code}/>
                                <InfoRow label="Country" value={(data.details.Address as Address)?.country}/>
                            </CardContent>
                        </Card>
                    </Box>
                )}
                {edit && (
                    <Box
                        sx={{
                            display: "grid",
                            gap: 3,
                            gridTemplateColumns: {xs: "1fr", md: "1fr 1fr"},
                        }}
                    >
                        <Card>
                            <CardContent>
                                <Typography variant="subtitle1" fontWeight={600} mb={2}>
                                    About
                                </Typography>
                                <Box
                                    sx={{
                                        display: "grid",
                                        gap: 3,
                                    }}
                                >
                                    <TextField
                                        size="small"
                                        label="First Name"
                                        value={form.firstName}
                                        onChange={(e) =>
                                            setForm((f) => ({...f, firstName: e.target.value}))
                                        }
                                    />
                                    <TextField
                                        size="small"
                                        label="Last Name"
                                        value={form.lastName}
                                        onChange={(e) =>
                                            setForm((f) => ({...f, lastName: e.target.value}))
                                        }
                                    />
                                    <TextField
                                        size="small"
                                        label="Phone"
                                        value={form.phone}
                                        onChange={(e) =>
                                            setForm((f) => ({...f, phone: e.target.value}))
                                        }
                                    />

                                    <Stack direction="row" spacing={2}>
                                        <Button
                                            variant="contained"
                                            startIcon={<Save/>}
                                            onClick={handleSave}
                                            disabled={saving}
                                        >
                                            {saving ? "Saving…" : "Save"}
                                        </Button>
                                        <Button onClick={() => setEdit(false)}>Cancel</Button>
                                    </Stack>
                                </Box>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardContent>
                                <Typography variant="subtitle1" fontWeight={600} mb={2}>
                                    Address
                                </Typography>
                                <Box
                                    sx={{
                                        display: "grid",
                                        gap: 3,
                                    }}
                                >
                                    <TextField
                                        size="small"
                                        label="Street / Line"
                                        value={form.address.line}
                                        onChange={(e) =>
                                            setForm(f => ({
                                                ...f,
                                                address: {
                                                    ...f.address,
                                                    line: e.target.value,
                                                },
                                            }))
                                        }
                                    />
                                    <TextField
                                        size="small"
                                        label="City"
                                        value={form.address.city}
                                        onChange={(e) =>
                                            setForm(f => ({
                                                ...f,
                                                address: {
                                                    ...f.address,
                                                    city: e.target.value,
                                                },
                                            }))
                                        }
                                    />
                                    <TextField
                                        size="small"
                                        label="Postal Code"
                                        value={form.address.postal_code}
                                        onChange={(e) =>
                                            setForm(f => ({
                                                ...f,
                                                address: {
                                                    ...f.address,
                                                    postal_code: e.target.value,
                                                },
                                            }))
                                        }
                                    />
                                    <TextField
                                        size="small"
                                        label="Country"
                                        value={form.address.country}
                                        onChange={(e) =>
                                            setForm(f => ({
                                                ...f,
                                                address: {
                                                    ...f.address,
                                                    country: e.target.value,
                                                },
                                            }))
                                        }
                                    />
                                </Box>
                            </CardContent>
                        </Card>
                    </Box>

                )}

                {/* EXPERIENCE & EDUCATION */}
                <Box
                    sx={{
                        display: "grid",
                        gap: 3,
                        gridTemplateColumns: {xs: "1fr", md: "1fr 1fr"},
                    }}
                >
                    <Card>
                        <CardContent>
                            <Typography variant="subtitle1" fontWeight={600} mb={1}>
                                Experience
                            </Typography>
                            <Typography variant="body2">
                                Owner at Her Company Inc.
                            </Typography>
                        </CardContent>
                    </Card>
                    <Card>
                        <CardContent>
                            <Typography variant="subtitle1" fontWeight={600} mb={1}>
                                Education
                            </Typography>
                            <Typography variant="body2">
                                Master’s Degree in Oxford
                            </Typography>
                        </CardContent>
                    </Card>
                </Box>
            </Box>

            {/* ---- dialog crop ---- */}
            {cropOpen && cropSrc && (
                <CropDialog
                    src={cropSrc}
                    onComplete={handleCropComplete}
                    onSave={handleCropSave}
                    onCancel={() => setCropOpen(false)}
                />
            )}
        </Box>
    );
}

/* --- sub-componentă InfoRow (eticheta + valoare) --- */
function InfoRow({label, value}: { label: string; value?: string }) {
    return (
        <Box
            sx={{
                display: "grid",
                gridTemplateColumns: "130px 1fr",
                columnGap: 2,
                mb: 1,
            }}
        >
            <Typography
                variant="body2"
                color="text.secondary"
                fontWeight={500}
                sx={{whiteSpace: "nowrap"}}
            >
                {label}
            </Typography>

            <Typography variant="body2">
                {value || "—"}
            </Typography>
        </Box>
    );
}
