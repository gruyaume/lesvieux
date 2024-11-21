"use client"
import { useRouter } from "next/navigation"

export default function FrontPage() {
    const router = useRouter()
    router.push("/employer_portal/login")
}