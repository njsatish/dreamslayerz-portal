import { useEffect, useState } from 'react'
import SiteLayout from './components/SiteLayout'
import HomePage from './pages/HomePage'
import ServicesPage from './pages/ServicesPage'
import AboutPage from './pages/AboutPage'
import BookPage from './pages/BookPage'
import ContactPage from './pages/ContactPage'
import GalleryPage from './pages/GalleryPage'
import ProviderPage from './pages/ProviderPage'
const API_BASE=import.meta.env.VITE_API_BASE_URL||'http://localhost:8000'
function route(){const value=window.location.hash.replace(/^#\/?/,'').split('?')[0];if(value.startsWith('providers/'))return 'provider';return ['home','services','about','book','gallery','contact'].includes(value)?value:'home'}
export default function App(){
 const[config,setConfig]=useState(null),[page,setPage]=useState(route())
 useEffect(()=>{fetch(`${API_BASE}/api/config`).then(r=>r.ok?r.json():Promise.reject()).then(setConfig).catch(()=>{})},[])
 useEffect(()=>{const change=()=>{setPage(route());window.scrollTo(0,0)};window.addEventListener('hashchange',change);return()=>window.removeEventListener('hashchange',change)},[])
 const navigate=target=>{window.location.hash=`/${target}`}
 const openGallery=itemId=>{window.location.hash=itemId?`/gallery?image=${encodeURIComponent(itemId)}`:'/gallery'}
 const providerId=decodeURIComponent((window.location.hash.match(/^#?\/providers\/([^?]+)/)||[])[1]||'')
 const content={home:<HomePage config={config} navigate={navigate} openGallery={openGallery}/>,services:<ServicesPage config={config} navigate={navigate}/>,about:<AboutPage config={config} navigate={navigate}/>,book:<BookPage/>,contact:<ContactPage config={config} navigate={navigate}/>,gallery:<GalleryPage onBack={()=>navigate('home')}/>,provider:<ProviderPage config={config} providerId={providerId} navigate={navigate}/>,provider:<ProviderPage config={config} providerId={providerId} navigate={navigate}/>}[page]
 return <SiteLayout page={page} navigate={navigate} business={config?.business}>{content}</SiteLayout>
}
