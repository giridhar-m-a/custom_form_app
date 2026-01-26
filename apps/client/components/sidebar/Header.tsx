import { useMemo } from 'react'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator
} from '../ui/breadcrumb'
import { usePathname } from 'next/navigation'
import { ThemeSwitch } from '../theme/ThemeSwitch'
import { Fragment } from 'react'

const Header = () => {
  const pathname = usePathname()
  const path = useMemo(() => {
    return pathname.split('/').filter(item => item !== '')
  }, [pathname])
  return (
    <div className="flex items-center gap-2 justify-between w-full">
      <Breadcrumb>
        <BreadcrumbList>
          {path.map((item, i) => (
            <Fragment key={i}>
              {i < path.length - 1 && (
                <BreadcrumbItem className="hidden md:block text-xl font-semibold capitalize">
                  <BreadcrumbLink href={`/${path.slice(0, i + 1).join('/')}`}>
                    {item.replaceAll('-', ' ')}
                  </BreadcrumbLink>
                </BreadcrumbItem>
              )}
              {i < path.length - 1 && <BreadcrumbSeparator className="hidden md:block" />}
              {i === path.length - 1 && (
                <BreadcrumbItem>
                  <BreadcrumbPage className="text-xl font-semibold capitalize">
                    {item.replaceAll('-', ' ')}
                  </BreadcrumbPage>
                </BreadcrumbItem>
              )}
            </Fragment>
          ))}
        </BreadcrumbList>
      </Breadcrumb>
      <ThemeSwitch />
    </div>
  )
}

export default Header
