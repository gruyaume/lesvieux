"use client"
import { useRouter } from "next/navigation"

export default function FrontPage() {
    const router = useRouter()
    router.push("/admin_portal/login")
}