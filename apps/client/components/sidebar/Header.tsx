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
import { useGetFormById } from '@/hooks/queryHooks/useFormApp'

const Header = () => {
  const pathname = usePathname()
  const path = useMemo(() => {
    return pathname.split('/').filter(item => item !== '')
  }, [pathname])
  const formId = path[0] === 'forms' && path[1] && path[1] !== 'new' && path[1] !== 'edit' ? path[1] : null
  const { data: formResponse } = useGetFormById(formId || '', { enabled: Boolean(formId) })
  const pageLabel = formResponse?.data?.title || path[path.length - 1]?.replaceAll('-', ' ')
  return (
    <div className="flex items-center gap-2 justify-between w-full">
      <Breadcrumb>
        <BreadcrumbList>
          {path.map((item, i) => (
            <Fragment key={i}>
              {i < path.length - 1 && (
                <BreadcrumbItem className="hidden md:block text-xl font-semibold capitalize">
                  {item !== 'edit' && item !== 'new' ? (
                    <BreadcrumbLink href={`/${path.slice(0, i + 1).join('/')}`}>
                      {item.replaceAll('-', ' ')}
                    </BreadcrumbLink>
                  ) : (
                    <BreadcrumbPage>{item.replaceAll('-', ' ')}</BreadcrumbPage>
                  )}
                </BreadcrumbItem>
              )}
              {i < path.length - 1 && <BreadcrumbSeparator className="hidden md:block" />}
              {i === path.length - 1 && (
                <BreadcrumbItem>
                  <BreadcrumbPage className="text-xl font-semibold capitalize">
                    {formId && i === path.length - 1 ? pageLabel : item.replaceAll('-', ' ')}
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
