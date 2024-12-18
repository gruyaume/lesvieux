"use client"

import { SetStateAction, Dispatch, useState, useContext } from "react"
import { useQuery } from "react-query";
import { AsideContext } from "./aside";
import { usePathname } from "next/navigation";
import { useAuth } from "./auth/authContext";
import { ChangePasswordModalData, ChangeMyPasswordModal, ChangePasswordModalContext } from "./components";
import Logo from "../components/logo"
import { Button, Panel, SideNavigation, StatusLabel } from "@canonical/react-components";
import { useCookies } from "react-cookie";
import { getStatus } from "../queries"


export function SideBar({ activePath, sidebarVisible }: { activePath: string, sidebarVisible: boolean, setSidebarVisible: Dispatch<SetStateAction<boolean>> }) {
    const changePasswordModalContext = useContext(ChangePasswordModalContext)
    const [cookies, setCookie, removeCookie] = useCookies(['user_token']);
    const [menuOpen, setMenuOpen] = useState<boolean>(false)
    const auth = useAuth()

    const parseVersion = (version: string) => {
        const [baseVersion, suffix] = version.split('-');
        return { baseVersion, suffix };
    };

    const statusQuery = useQuery({
        queryFn: getStatus,
        staleTime: Infinity,
        cacheTime: Infinity,
        refetchOnWindowFocus: false,
        refetchOnMount: false,
        refetchOnReconnect: false,
    })

    if (!auth.user) {
        return <></>;
    }

    const versionInfo = parseVersion(statusQuery.data?.version || "");

    return (
        <header className={sidebarVisible ? "l-navigation" : "l-navigation"}>
            <Panel
                stickyHeader
                title={
                    <div >
                        <Logo />
                        <StatusLabel
                            appearance="information">
                            Employer
                        </StatusLabel>
                    </div>
                }
            >
                <SideNavigation
                    hasIcons
                    items={[
                        {
                            items: [
                                {
                                    href: '/employer_portal/my_posts',
                                    "aria-current": activePath.startsWith("/employer_portal/my_posts"),
                                    icon: 'canvas',
                                    label: 'Job Posts'
                                },
                            ]
                        }
                    ]}
                />

                <div className="p-side-navigation--icons" id="drawer-icons">
                    <nav aria-label="Main">
                        <ul className="p-side-navigation__list" style={{ bottom: "64px", position: "absolute", width: "100%" }}>
                            <li className="p-side-navigation__item" >
                                <div className="p-side-navigation__link p-contextual-menu__toggle" onClick={() => setMenuOpen(!menuOpen)} aria-current={menuOpen} style={{ cursor: "pointer" }}>
                                    <i className="p-icon--user p-side-navigation__icon"></i>
                                    <span className="p-side-navigation__label">
                                        <span className="p-side-navigation__label">{auth.user.email}</span>
                                    </span>
                                    <div className="p-side-navigation__status">
                                        <i className="p-icon--menu"></i>
                                        <span className="p-contextual-menu__dropdown" id="menu-3" aria-hidden={!menuOpen} style={{ bottom: "40px" }}>
                                            <span className="p-contextual-menu__group">
                                                <Button
                                                    className="p-contextual-menu__link"
                                                    onMouseDown={() => changePasswordModalContext.setModalData({ "id": auth.user ? auth.user.id.toString() : "", "email": auth.user ? auth.user.email : "" })}>
                                                    Change Password
                                                </Button>
                                                <Button
                                                    className="p-contextual-menu__link"
                                                    onMouseDown={() => removeCookie("user_token")}>
                                                    Log Out
                                                </Button>
                                            </span>
                                        </span>
                                    </div>
                                </div>
                            </li>
                        </ul>
                        <ul className="p-side-navigation__list" style={{ bottom: 0, position: "absolute", width: "100%" }}>
                            <li className="p-side-navigation__item">
                                <span className="p-side-navigation__text">
                                    {`Version ${versionInfo.baseVersion}`}
                                    {versionInfo.suffix && (
                                        <div className="p-side-navigation__status">
                                            <StatusLabel
                                                appearance="caution">
                                                {versionInfo.suffix}
                                            </StatusLabel>
                                        </div>
                                    )}
                                </span>
                            </li>
                        </ul>
                    </nav>
                </div>
            </Panel>
        </header >
    )
}


export default function Navigation({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    const activePath = usePathname()
    const noNavRoutes = ['/employer_portal/login', '/employer_portal/my_posts/draft'];

    const shouldRenderNavigation = !noNavRoutes.includes(activePath);
    const [sidebarVisible, setSidebarVisible] = useState<boolean>(true)
    const [asideOpen, setAsideOpen] = useState<boolean>(false)
    const [asideData, setAsideData] = useState<any>(null)
    const [changePasswordModalData, setChangePasswordModalData] = useState<ChangePasswordModalData>(null)

    return (
        <div className="l-application" role="presentation">
            <AsideContext.Provider value={{ isOpen: asideOpen, setIsOpen: setAsideOpen, extraData: asideData, setExtraData: setAsideData }}>
                <ChangePasswordModalContext.Provider value={{ modalData: changePasswordModalData, setModalData: setChangePasswordModalData }}>
                    {
                        shouldRenderNavigation ? (
                            <>
                                <SideBar activePath={activePath} sidebarVisible={sidebarVisible} setSidebarVisible={setSidebarVisible} />
                            </>
                        ) : (
                            <></>
                        )
                    }
                </ChangePasswordModalContext.Provider>
                <main className="l-main">
                    {children}
                    {changePasswordModalData != null && <ChangeMyPasswordModal modalData={changePasswordModalData} setModalData={setChangePasswordModalData} />}
                </main>
            </AsideContext.Provider>
        </div >
    )
}