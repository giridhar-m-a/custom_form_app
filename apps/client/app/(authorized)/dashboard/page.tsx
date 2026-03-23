'use client'
import FileUpload from '@/components/FormRender/FileUpload'

const Dashboard = () => {
  const handlePath = (path: string) => {
    console.log('path: ', path)
  }
  return (
    <div>
      <h1>Dashboard</h1>
      <FileUpload uploadPath="abc/abc" autoUpload={true} handlePath={handlePath} />
    </div>
  )
}

export default Dashboard
