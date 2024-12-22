import { useTranslation } from 'react-i18next'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { buttonVariants } from '@/components/ui/button'
import { HomeIcon } from 'lucide-react';

function NotFound() {
  const { t } = useTranslation();

  return (
    <>
      <div className="flex justify-center items-center h-[calc(100vh-4rem)]">
        <Card className="max-w-[400px] w-full mx-4">
          <CardHeader>
            <CardTitle>{t('player.404.title')}</CardTitle>
          </CardHeader>
          <CardContent>
            <CardDescription dangerouslySetInnerHTML={{ __html: t('player.404.description') }} />
            <div className="flex justify-end items-center w-full">
              <a href="/" className={buttonVariants({ variant: "outline" }) + ` mt-4`}>
                <HomeIcon className="w-4 h-4 mr-1" />
                {t('player.404.back')}
              </a>
            </div>
          </CardContent>
        </Card>
      </div>
    </>
  )
}

export default NotFound
