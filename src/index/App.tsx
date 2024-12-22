import moment from 'moment'
import { useEffect, useLayoutEffect, useState } from 'react'
import { t } from 'i18next'
import {
  HomeIcon,
  FileIcon,
  FileTextIcon,
  FileImageIcon,
  FileAudio2Icon,
  FileVideo2Icon,
  FileCodeIcon,
  FileJsonIcon,
  FileArchiveIcon,
  FileDigitIcon,
  FolderIcon,
  ArrowDownAzIcon,
  ArrowUpAzIcon,
  ArrowDown01Icon,
  ArrowUp01Icon,
  ClockArrowDownIcon,
  ClockArrowUpIcon,
  MonitorPlayIcon,
  ChevronRightIcon,
  SlidersHorizontalIcon,
  BanIcon,
} from 'lucide-react'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
} from '@/components/ui/breadcrumb'
import {
  Card,
  CardContent,
  CardHeader,
} from '@/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger
} from '@/components/ui/context-menu'
import {
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select } from '@/components/ui/select'
import { Button } from '@/components/ui/button'
import { DirectoryInfo, FileInfo, FileType, GetFileType } from '@/lib/directory'
import Footer from '@/components/footer'
import { encodeURIRFC3986 } from  '@/lib/utils'
import './App.css'


type TimezoneType = 'client' | 'server' | 'utc'
interface SortState {
  key: 'name' | 'size' | 'lastModified',
  order: 'asc' | 'desc'
}
const DEFAULT_TITLE = import.meta.env.APP_NAME


function App() {
  const FILEICON_MAP: Record<FileType, React.ReactNode> = {
    'folder': <FolderIcon className="w-5 h-5" />,
    'document': <FileTextIcon className="w-5 h-5" />,
    'image': <FileImageIcon className="w-5 h-5" />,
    'audio': <FileAudio2Icon className="w-5 h-5" />,
    'video': <FileVideo2Icon className="w-5 h-5" />,
    'code': <FileCodeIcon className="w-5 h-5" />,
    'config': <FileJsonIcon className="w-5 h-5" />,
    'archive': <FileArchiveIcon className="w-5 h-5" />,
    'binary': <FileDigitIcon className="w-5 h-5" />,
    'file': <FileIcon className="w-5 h-5" />,
  }

  const [directoryInfo, setDirectoryInfo] = useState<DirectoryInfo | null>(null)
  const [sortState, setSortState] = useState<SortState>({ key: 'name', order: 'asc' })
  const [timezone, setTimezone] = useState<TimezoneType>('client')
  const [search, setSearch] = useState<string>('')

  function XMLQuerySelector(xml: Document | Element, selector: string, defaultValue: string) {
    const element = xml.querySelector(selector) || xml.querySelector(selector.toLowerCase()) || xml.querySelector(selector.toUpperCase())
    if (!element) return defaultValue
    return element.textContent || defaultValue
  }

  function XMLQuerySelectorAll(xml: Document, selector: string) {
    let elements = xml.querySelectorAll(selector)
    if (!elements.length) {
      elements = xml.querySelectorAll(selector.toLowerCase())
    }
    if (!elements.length) {
      elements = xml.querySelectorAll(selector.toUpperCase())
    }
    return elements
  }

  useEffect(() => {
    // Firefox does not like textContent. Make it happy.
    const data = document.querySelector('#xml-data')?.outerHTML
    if (!data) return
    const parser = new DOMParser()
    // DOMParser does not select root element first like XPath, so one code for all browsers, happy?
    const xml = parser.parseFromString(data, 'application/xml')
    const info: DirectoryInfo = {
      bucketName: XMLQuerySelector(xml, 'Name', DEFAULT_TITLE),
      path: XMLQuerySelector(xml, 'Prefix', '/'),
      serverTimezoneOffset: XMLQuerySelector(xml, 'ServerTimezoneOffset', '+00:00'),
      files: Array.from(XMLQuerySelectorAll(xml, 'Contents')).map(file => ({
        isDirectory: XMLQuerySelector(file, 'IsDirectory', 'false') === 'true',
        fileName: XMLQuerySelector(file, 'FileName', ''),
        lastModifiedUnix: parseInt(XMLQuerySelector(file, 'LastModifiedUnix', '0')),
        size: parseInt(XMLQuerySelector(file, 'Size', '0')),
      }))
    }
    setDirectoryInfo(info)
    setTimezone('client')
  }, [])

  useLayoutEffect(() => {
    // Fix chrome XML render bug
    if (!(window as any).chrome) return
    const windowTimer = setInterval(() => {
      const html = document.querySelector("html")
      if (!html) return
      console.log(html.clientHeight, html.scrollHeight, "client >= scroll", html.clientHeight >= html.scrollHeight)
      if (html.clientHeight < html.scrollHeight) {
        html.style.height = html.scrollHeight + 20 + "px"
      }
      if (html.clientHeight >= html.scrollHeight) {
        console.log('clearInterval')
        clearInterval(windowTimer)
      }
    }, 100)
    return () => clearInterval(windowTimer)
  }, [])

  useEffect(() => {
    document.title = directoryInfo?.bucketName || DEFAULT_TITLE
  }, [directoryInfo])

  function getHumanReadableTime(time: number) {
    let offset = '+00:00'
    if (timezone == 'client') {
      offset = moment().format('Z')
    }
    if (timezone == 'server') {
      offset = directoryInfo?.serverTimezoneOffset || offset
    }
    return moment.unix(time).utcOffset(offset).format('YYYY-MM-DD HH:mm:ss')
  }

  function getHumanReadableSize(size: number) {
    const units = ['B', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB']
    let i = 0
    while (size >= 1024) {
      size /= 1024
      i++
    }
    return size.toFixed(2) + ' ' + units[i]
  }

  function filterFiles() {
    return directoryInfo?.files.filter(file => file.fileName.toLowerCase().includes(search.toLowerCase())) || []
  }

  function sortedSearchingFiles(): FileInfo[] {
    const files = filterFiles()

    if (sortState.key === 'name') {
      return files.sort((a, b) => {
        // 1. Put folders first
        if (a.isDirectory && !b.isDirectory) return -1
        if (!a.isDirectory && b.isDirectory) return 1
        // 2. Sort by file name
        const sign = sortState.order === 'asc' ? 1 : -1
        const nameCompare = a.fileName.localeCompare(b.fileName)
        if (nameCompare !== 0) return sign * nameCompare
        // 3. If file name is the same, sort by last modified time
        if (a.lastModifiedUnix !== b.lastModifiedUnix) {
          return a.lastModifiedUnix - b.lastModifiedUnix
        }
        // 4. If last modified time is the same, sort by file size
        return a.size - b.size
      })
    }

    if (sortState.key === 'lastModified') {
      return files.sort((a, b) => {
        // 1. Put folders first
        if (a.isDirectory && !b.isDirectory) return -1
        if (!a.isDirectory && b.isDirectory) return 1
        // 2. Sort by last modified time
        const sign = sortState.order === 'asc' ? 1 : -1
        const timeCompare = a.lastModifiedUnix - b.lastModifiedUnix
        if (timeCompare !== 0) return sign * timeCompare
        // 3. If last modified time is the same, sort by file name
        const nameCompare = a.fileName.localeCompare(b.fileName)
        if (nameCompare !== 0) return nameCompare
        // 4. If file name is the same, sort by file size
        return a.size - b.size
      })
    }

    if (sortState.key === 'size') {
      return files.sort((a, b) => {
        // 1. Put folders first
        const sign = sortState.order === 'asc' ? 1 : -1
        if (a.isDirectory && !b.isDirectory) return -1 * sign
        if (!a.isDirectory && b.isDirectory) return 1 * sign
        // 2. Sort by file size
        const sizeCompare = a.size - b.size
        if (sizeCompare !== 0) return sign * sizeCompare
        // 3. If file size is the same, sort by file name
        const nameCompare = a.fileName.localeCompare(b.fileName)
        if (nameCompare !== 0) return nameCompare
        // 4. If file name is the same, sort by last modified time
        return a.lastModifiedUnix - b.lastModifiedUnix
      })
    }

    // This should never happen, if it does, return the original order
    return files
  }

  function toggleSort(key: 'name' | 'size' | 'lastModified') {
    if (sortState.key === key) {
      setSortState({ key, order: sortState.order === 'asc' ? 'desc' : 'asc' })
    } else {
      setSortState({ key, order: 'asc' })
    }
  }

  function playUrl(filename: string) {
    const filePath = `${directoryInfo?.path.replace(/\/$/, '')}/${encodeURIRFC3986(filename)}`.replace(/^\//, '')
    return `/?_/player/${filePath}`
  }

  return (
    <>
      <header className="px-8 py-4 bg-background w-full border-b shadow-md">
        <nav className="mt-2 mb-4 flex items-center justify-between">
          <h1 className="md:text-2xl text-xl">
            {directoryInfo?.bucketName || DEFAULT_TITLE}
          </h1>
          <Popover>
            <PopoverTrigger asChild>
              <Button variant="ghost" className="p-2">
                <SlidersHorizontalIcon className="w-6 h-6" />
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-80" align="end">
              <div className="w-full">
              <Label className="block ml-1 mb-2">{t('index.search')}</Label>
              <Input className="m-0 flex-grow" placeholder={t('index.search_file') + '...'} name={t('index.search_file')} value={search} onChange={e => setSearch(e.target.value)} />
            </div>
            <div className="w-full mt-6">
              <Label className="block ml-1 mb-2">{t('index.timezone')}</Label>
              <Select value={timezone} onValueChange={value => setTimezone(value as TimezoneType)}>
                <SelectTrigger className="w-full">
                  <SelectValue placeholder={t('index.timezone')} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="client">{t('index.timezone_client')} ({moment().format('Z')})</SelectItem>
                  <SelectItem value="server">{t('index.timezone_server')} ({directoryInfo?.serverTimezoneOffset})</SelectItem>
                  <SelectItem value="utc">{t('index.timezone_server')}</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="w-full mt-6 mb-2">
              <Label className="block ml-1 mb-2">{t('index.sort')}</Label>
              { sortState.key === 'name' && sortState.order === 'asc' && (
                <Button className="w-full" variant="outline" onClick={() => setSortState({ key: 'name', order: 'desc' })}>
                  {t('index.filename')} <ArrowDownAzIcon className="w-4 h-4 ml-1" />
                </Button>
              )}
              { sortState.key === 'name' && sortState.order === 'desc' && (
                <Button className="w-full" variant="outline" onClick={() => setSortState({ key: 'size', order: 'asc' })}>
                  {t('index.filename')} <ArrowUpAzIcon className="w-4 h-4 ml-1" />
                </Button>
              )}
              { sortState.key === 'size' && sortState.order === 'asc' && (
                <Button className="w-full" variant="outline" onClick={() => setSortState({ key: 'size', order: 'desc' })}>
                  {t('index.size')} <ArrowDown01Icon className="w-4 h-4 ml-1" />
                </Button>
              )}
              { sortState.key === 'size' && sortState.order === 'desc' && (
                <Button className="w-full" variant="outline" onClick={() => setSortState({ key: 'lastModified', order: 'asc' })}>
                  {t('index.size')} <ArrowUp01Icon className="w-4 h-4 ml-1" />
                </Button>
              )}
              { sortState.key === 'lastModified' && sortState.order === 'asc' && (
                <Button className="w-full" variant="outline" onClick={() => setSortState({ key: 'lastModified', order: 'desc' })}>
                  {t('index.last_modified')} <ClockArrowUpIcon className="w-4 h-4 ml-1" />
                </Button>
              )}
              { sortState.key === 'lastModified' && sortState.order === 'desc' && (
                <Button className="w-full" variant="outline" onClick={() => setSortState({ key: 'name', order: 'asc' })}>
                  {t('index.last_modified')} <ClockArrowDownIcon className="w-4 h-4 ml-1" />
                </Button>
              )}
            </div>
            </PopoverContent>
          </Popover>
        </nav>
        <Breadcrumb className="mb-1">
          <BreadcrumbList className="whitespace-nowrap overflow-x-auto flex-nowrap breadcrumbs">
            <BreadcrumbItem>
              <BreadcrumbLink href="/">
                <HomeIcon className="w-4 h-4" />
              </BreadcrumbLink>
            </BreadcrumbItem>

            {directoryInfo?.path !== '/' && (
              <>
                {directoryInfo?.path.split('/').map((dir, index, dirs) => (
                  <BreadcrumbItem key={index}>
                    <ChevronRightIcon className="w-4 h-4" />
                    <BreadcrumbLink href={`/${dirs.slice(0, index + 1).join('/')}`} className="flex items-center gap-1">
                      <FolderIcon className="w-4 h-4" />
                      {dir}
                    </BreadcrumbLink>
                  </BreadcrumbItem>
                ))}
              </>
            )}
          </BreadcrumbList>
        </Breadcrumb>
      </header>

      <Card className="min-h-[50vh] max-w-[1200px] lg:w-[80%] w-[90%] md:mx-auto mt-8 border md:block hidden">
        <CardHeader className="flex flex-row items-center justify-between gap-4">
          <div className="w-full" style={{ marginTop: '.5rem' }}>
            <Label className="block ml-1 mb-2">{t('index.search')}</Label>
            <Input className="m-0 flex-grow" placeholder={t('index.search_file') + '...'} name={t('index.search_file')} value={search} onChange={e => setSearch(e.target.value)} />
          </div>
          <div style={{ marginTop: '.5rem' }}>
            <Label className="block ml-1 mb-2">{t('index.timezone')}</Label>
            <Select value={timezone} onValueChange={value => setTimezone(value as TimezoneType)}>
              <SelectTrigger className="w-48">
                <SelectValue placeholder={t('index.timezone')} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="client">{t('index.timezone_client')} ({moment().format('Z')})</SelectItem>
                <SelectItem value="server">{t('index.timezone_server')} ({directoryInfo?.serverTimezoneOffset})</SelectItem>
                <SelectItem value="utc">{t('index.timezone_utc')}</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardHeader>
        <CardContent>
          <Table className="table-fixed w-full border-collapse">
            <TableHeader>
              <TableRow className="w-full">
                <TableHead className="px-0 w-7" />
                <TableHead className="px-0">
                  <Button className="w-full force-justify-left" variant="ghost" onClick={() => toggleSort('name')}>
                    {t('index.filename')}
                    {sortState.key === 'name' && sortState.order === 'asc' && <ArrowDownAzIcon className="w-4 h-4 ml-1" />}
                    {sortState.key === 'name' && sortState.order === 'desc' && <ArrowUpAzIcon className="w-4 h-4 ml-1" />}
                  </Button>
                </TableHead>
                <TableHead className="px-0 w-28">
                  <Button className="w-full force-justify-right" variant="ghost" onClick={() => toggleSort('size')}>
                    {sortState.key === 'size' && sortState.order === 'asc' && <ArrowDown01Icon className="w-4 h-4 ml-1" />}
                    {sortState.key === 'size' && sortState.order === 'desc' && <ArrowUp01Icon className="w-4 h-4 ml-1" />}
                    {t('index.size')}
                  </Button>
                </TableHead>
                <TableHead className="px-0 w-48">
                  <Button className="w-full force-justify-right" variant="ghost" onClick={() => toggleSort('lastModified')}>
                    {sortState.key === 'lastModified' && sortState.order === 'asc' && <ClockArrowUpIcon className="w-4 h-4 ml-1" />}
                    {sortState.key === 'lastModified' && sortState.order === 'desc' && <ClockArrowDownIcon className="w-4 h-4 ml-1" />}
                    {t('index.last_modified')}
                  </Button>
                </TableHead>
                <TableHead className="px-0 w-7" />
              </TableRow>
            </TableHeader>
            <TableBody>
              {sortedSearchingFiles().map((file, index) => (
                <TableRow key={index} className="w-full">
                  <TableCell className="w-7 pl-0.5">
                    {FILEICON_MAP[GetFileType(file)]}
                  </TableCell>
                  <TableCell className="truncate p-0">
                    <a href={encodeURIRFC3986(file.fileName)}
                      className="truncate block pl-2 w-full h-8 flex items-center font-medium hover:text-blue-500">
                      <span className="truncate">{file.fileName}</span>
                    </a>
                  </TableCell>
                  <TableCell className="whitespace-nowrap text-right w-28 code">{file.isDirectory ? '-' : getHumanReadableSize(file.size)}</TableCell>
                  <TableCell className="whitespace-nowrap text-right w-48 code">{getHumanReadableTime(file.lastModifiedUnix)}</TableCell>
                  <TableCell className="w-7">
                    {GetFileType(file) === 'video' && (
                      <a href={playUrl(file.fileName)}
                        className="hover:text-blue-500">
                        <MonitorPlayIcon className="w-5 h-5" />
                      </a>
                    )}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card >

      <div className="md:hidden block">
        {sortedSearchingFiles().map((file, index) => (
          <ContextMenu key={index}>
            <ContextMenuTrigger>
              <a href={encodeURIRFC3986(file.fileName)} className="w-full flex p-4 items-center border-b-2 cursor-pointer">
              <div>
                {FILEICON_MAP[GetFileType(file)]}
              </div>
              <div className="w-full flex flex-col ml-4 truncate">
                <div className="w-full truncate mb-0.5">{file.fileName}</div>
                <div className="w-full flex items-center justify-between">
                  <div className="text-xs text-gray-400 code">{file.isDirectory ? '-' : getHumanReadableSize(file.size)}</div>
                  <div className="text-xs text-gray-400 code">{getHumanReadableTime(file.lastModifiedUnix)}</div>
                </div>
              </div>
            </a>
            </ContextMenuTrigger>
            <ContextMenuContent>
              {GetFileType(file) === 'video' && (
                <ContextMenuItem className="flex">
                  <a href={playUrl(file.fileName)} className="flex items-center text-lg  gap-2">
                    <MonitorPlayIcon />
                    Play
                  </a>
                </ContextMenuItem>
              )}
              <ContextMenuItem className="flex items-center text-lg gap-2">
                <BanIcon />
                Cancel
              </ContextMenuItem>
            </ContextMenuContent>
          </ContextMenu>
        ))}
      </div>

      <Footer />
    </>
  )
}

export default App
