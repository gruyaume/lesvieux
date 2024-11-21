import { SetStateAction, Dispatch, createContext, useContext } from "react"
import { useAuth } from "./auth/authContext"
import { AppAside } from "@canonical/react-components";

type AsideContextType = {
    isOpen: boolean,
    setIsOpen: Dispatch<SetStateAction<boolean>>

    extraData: any
    setExtraData: Dispatch<SetStateAction<any>>
}

export const AsideContext = createContext<AsideContextType>({
    isOpen: false,
    setIsOpen: () => { },

    extraData: null,
    setExtraData: () => { },
})

export function Aside({ FormComponent }: { FormComponent: React.ComponentType<any> }) {
    const auth = useAuth()
    const asideContext = useContext(AsideContext)
    return (
        <AppAside
            className={(auth.user && asideContext.isOpen ? "" : " is-collapsed")}
        >
            <FormComponent />
        </AppAside >
    )
}