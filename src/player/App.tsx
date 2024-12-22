import { useEffect, useLayoutEffect, useRef, useState } from 'react'
import { useToast } from '@/hooks/use-toast'
import ArtPlayer from 'artplayer'
import ArtPlayerPluginDanmuku from 'artplayer-plugin-danmuku'
import {
  AirplayIcon, ArrowRightToLineIcon, CameraIcon, CaptionsIcon, CheckIcon, DownloadIcon,
  FolderTreeIcon, KeyboardIcon, KeyboardOffIcon, PictureInPictureIcon, RepeatIcon,
  SlidersHorizontalIcon,
} from 'lucide-react'
import { Card, CardContent } from '@/components/ui/card'
import { AspectRatio } from '@/components/ui/aspect-ratio'
import { Toaster } from '@/components/ui/toaster'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import Footer from '@/components/footer'
import NotFound from './NotFound'
import { Button, buttonVariants } from '@/components/ui/button'
import { t } from 'i18next'
import './App.css'


const PLAYER_PREFIX = '?_/player'
const DESKTOP_MODE_MIN_WIDTH = 840
const AUTOPLAY_WARN_TIMEOUT_MS = 3000
const APP_NAME = import.meta.env.APP_NAME


interface SubtitleData {
  langName?: string
  langInfo?: any
  lang: string
  url: string
}

interface PlayerData {
  danmuku: string | null
  subtitles: SubtitleData[]
}


function App() {
  const { toast } = useToast()
  const artRef = useRef<HTMLDivElement>(null)
  const danmukuRef = useRef<HTMLDivElement>(null)
  const [player, setPlayer] = useState<ArtPlayer | null>(null)
  const [videoURL, setVideoURL] = useState('#')
  const [is404, setIs404] = useState(false)
  const [mobileMode, setMobileMode] = useState(false)
  const [isDanmukuVisible, setIsDanmukuVisible] = useState(true)
  const [videoLoop, setVideoLoopUnderleyer] = useState(true)
  const [videoSubtitle, setVideoSubtitleUnderleyer] = useState<SubtitleData | null>(null)

  function getMetchedLenguageSubtitles(subtitles: SubtitleData[]): SubtitleData | null {
    const browserLang = navigator.language.replace(/-/, '_')
    const langFormats = [browserLang, browserLang.replace(/_/, '-'), browserLang.replace(/_/, ''), browserLang.split('_')[0]]
    for (const lang of langFormats) {
      for (const subtitle of subtitles) {
        if (subtitle.lang.toLowerCase().startsWith(lang.toLowerCase())) {
          return subtitle
        }
      }
    }
    return null
  }

  function setSubtitle(subtitle: SubtitleData | null) {
    console.log('Set subtitle', subtitle, player?.subtitle)
    setVideoSubtitleUnderleyer(subtitle)
    if (!player) { return }
    player.subtitle.switch(subtitle?.url || '')
  }

  function setVideoLoop(loop: boolean) {
    console.log('Set video loop', loop, player?.option, player?.video)
    setVideoLoopUnderleyer(loop)
    if (player) {
      player.option.loop = loop
      player.video.loop = loop
    }
  }

  function isMoblieMode() {
    // Shared by React and ArtPlayer
    return window.innerWidth < DESKTOP_MODE_MIN_WIDTH
  }

  function updateMobileMode() {
    setMobileMode(isMoblieMode())
  }

  function toggleDanmuku() {
    const danmukuPlugin = player?.plugins.artplayerPluginDanmuku
    if (!danmukuPlugin) { return }
    if (danmukuPlugin.isHide) {
      danmukuPlugin.show()
      setIsDanmukuVisible(true)
    } else {
      danmukuPlugin.hide()
      setIsDanmukuVisible(false)
    }
  }

  useLayoutEffect(() => {
    updateMobileMode()
    window.addEventListener('resize', updateMobileMode)
    return () => {
      window.removeEventListener('resize', updateMobileMode)
    }
  }, [])

  function getPlayerData(): PlayerData | null {
    const result = {
      danmuku: null,
      subtitles: [],
    } as PlayerData
    const playerData = document.querySelector('#player-data')?.textContent
    if (!playerData) { return null }
    try {
      const parsed = JSON.parse(playerData) as any
      if (parsed.danmaku) { result.danmuku = parsed.danmaku }
      if (parsed.danmuku) { result.danmuku = parsed.danmuku }  // Handle multiple-translated versions
      if (parsed.subtitles && Array.isArray(parsed.subtitles)) {
        parsed.subtitles.forEach((subtitle: any) => {
          if (subtitle.lang && subtitle.url && typeof subtitle.lang === 'string' && typeof subtitle.url === 'string') {
            result.subtitles.push({
              lang: subtitle.lang,
              url: subtitle.url,
              langName: subtitle.lang_name,
              langInfo: subtitle.lang_info,
            })
          }
        })
      }
      return result
    } catch (e) { return null }
  }

  function getVideoFolderURL() {
    return videoURL.split('/').slice(0, -1).join('/') || '/'
  }

  useEffect(() => {
    const url = new URL(window.location.href)
    if (!url.search.startsWith(PLAYER_PREFIX)) {
      setIs404(true)
      return
    }

    const parsedVideoURL = url.search.slice(PLAYER_PREFIX.length)
    console.log('Video URL', parsedVideoURL)
    if (!parsedVideoURL) { return setIs404(true) }
    setVideoURL(parsedVideoURL)


    const playerInfo = getPlayerData()
    const matchingSubtitle = getMetchedLenguageSubtitles(playerInfo?.subtitles || [])
    console.log('Matching subtitle is', matchingSubtitle)

    ArtPlayer.PLAYBACK_RATE = [0.5, 0.8, 1, 1.5, 2, 2.5, 3, 4]
    const art = new ArtPlayer({
      container: artRef.current!,
      url: parsedVideoURL,
      autoplay: true,
      autoOrientation: true,
      setting: true,
      flip: true,
      loop: false,
      fullscreen: true,
      fullscreenWeb: true,
      playbackRate: true,
      aspectRatio: true,
      subtitleOffset: true,
      fastForward: true,
      plugins: [
        ArtPlayerPluginDanmuku({
          danmuku: playerInfo?.danmuku || '',
          emitter: false,
          heatmap: true,
          fontSize: '5%',
          synchronousPlayback: true,
          mount: danmukuRef.current!,
        }),
      ],
    })
    art.video.loop = videoLoop
    art.subtitle.switch(matchingSubtitle?.url || '')
    setSubtitle(matchingSubtitle)

    art.on('resize', () => {
      if (isMoblieMode()) {
        art.plugins.artplayerPluginDanmuku.config({ fontSize: '5%' })
      } else {
        art.plugins.artplayerPluginDanmuku.config({ fontSize: 25 })
      }
    })

    const autoplayWarnTimeout = setTimeout(() => {
      if (!art.isReady) { return; }
      if (!art.playing && art.played == 0) {
        console.error('Autoplay failed, show warning')
        toast({
          title: t('player.autoplay_warn.title'),
          description: t('player.autoplay_warn.description'),
        })
        return
      }
      console.log('ArtPlay state', art.playing, art.played)
    }, AUTOPLAY_WARN_TIMEOUT_MS)
    setPlayer(art)

    return () => {
      if (art && art.destroy) {
        art.destroy(false)
      }
      clearTimeout(autoplayWarnTimeout)
    }
  }, [])

  return (
    <>
      <Toaster />

      {is404 ? <NotFound /> :
        <>
          <Card className={mobileMode ? 'mobile' : 'desktop'}>
            <CardContent>
              <div className="flex mb-4 items-center justify-between player-banner-desktop">
                <a href={getVideoFolderURL()} className={buttonVariants({ variant: "outline" }) + ` mt-6`}>
                  <FolderTreeIcon />
                  <span>{t('player.view_folder')}</span>
                </a>
                <h1 className="text-xl flex items-center mt-6">
                  {APP_NAME} {t('player.title')}
                </h1>
                <a href={videoURL} target="_blank" download
                  className={buttonVariants({ variant: "outline" }) + ` mt-6`}>
                  <DownloadIcon />
                  <span>{t('player.download_video')}</span>
                </a>
              </div>
              <AspectRatio ratio={16 / 9} ref={artRef} className="artplayer-app w-full grow-0 flex m-0" />
              <div ref={danmukuRef} className="hidden" />
              <div className="mt-4 flex items-center justify-between mx-2 gap-4 buttom-controls">
                <div className="flex gap-4 mobile-only-controls">
                  <a href={getVideoFolderURL()} className={buttonVariants({ variant: "outline" })}>
                    <FolderTreeIcon />
                  </a>
                  <a href={videoURL} target="_blank" download
                    className={buttonVariants({ variant: "outline" })}>
                    <DownloadIcon />
                  </a>
                </div>
                <div className="flex gap-4 desktop-only-controls" />
                <div className="flex gap-4">
                <Popover>
                    <PopoverTrigger asChild>
                      <Button variant="outline" className="p-2">
                        <CaptionsIcon />
                        {t('player.subtitles')}
                      </Button>
                    </PopoverTrigger>
                    <PopoverContent className="w-40 flex flex-col p-1" align="end" style={{ zIndex: 100 }}>
                      {getPlayerData()?.subtitles.map((subtitle, index) => (
                        <Button key={index} variant="ghost" onClick={() => setSubtitle(subtitle)} className="flex items-center space-between text-sm p-1">
                          {videoSubtitle?.lang === subtitle.lang ? <CheckIcon className='w-4' /> : <span className="w-4" />}
                          <span className="w-[60%] block truncate">{subtitle.langName || subtitle.lang}</span>
                        </Button>
                      ))}
                    </PopoverContent>
                  </Popover>
                  <Popover>
                    <PopoverTrigger asChild>
                      <Button variant="outline" className="p-2">
                        <SlidersHorizontalIcon />
                        {t('player.settings')}
                      </Button>
                    </PopoverTrigger>
                    <PopoverContent className="w-60 flex flex-col gap-4" align="end" style={{ zIndex: 100 }}>
                      <Button variant="outline" onClick={toggleDanmuku}>
                        {isDanmukuVisible ? <KeyboardIcon /> : <KeyboardOffIcon />}
                        <span className="w-[60%] block">{t('player.danmuku')}</span>
                      </Button>
                      <Button variant="outline" onClick={() => setVideoLoop(!videoLoop)}>
                        { videoLoop ? <RepeatIcon /> : <ArrowRightToLineIcon /> }
                        <span className="w-[60%] block">{t('player.loop')}</span>
                      </Button>
                      <Button variant="outline" onClick={async () => player?.screenshot()}>
                        <CameraIcon />
                        <span className="w-[60%] block">{t('player.screenshot')}</span>
                      </Button>
                      <Button variant="outline" onClick={() => player?.airplay()}>
                        <AirplayIcon />
                        <span className="w-[60%] block">{t('player.airplay')}</span>
                      </Button>
                      <Button variant="outline" onClick={() => { player && (player.pip = !player.pip) }}>
                        <PictureInPictureIcon />
                        <span className="w-[60%] block">{t('player.pip')}</span>
                      </Button>
                    </PopoverContent>
                  </Popover>
                </div>
              </div>
            </CardContent>
          </Card>
        </>
      }

      <Footer />
    </>
  )
}

export default App
