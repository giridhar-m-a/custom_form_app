'use client'

import { useUploadProfilePic } from '@/hooks/queryHooks/useUsers'
import { getFileUrl, validateFile } from '@/lib/utils'
import { User } from '@/types/user.types'
import { Camera, User as UserIcon } from 'lucide-react'
import Image from 'next/image'
import { ChangeEvent, useRef, useState } from 'react'
import toast from 'react-hot-toast'

export const ProfilePic = ({ user }: { user?: User }) => {
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [image, setImage] = useState<File | null>(null)

  const handleImageClick = () => {
    fileInputRef.current?.click()
  }

  const { mutate, isPending } = useUploadProfilePic()

  const handleImageChange = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    const { error, valid } = validateFile({
      file,
      size: 10000000,
      type: new Set(['image/jpeg', 'image/png', 'image/webp', 'image/jpg']),
      setFile: setImage
    })
    if (!valid && error) {
      toast.error(error)
      return
    }
    const formData = new FormData()
    formData.append('file', file!)
    mutate(formData)
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  return (
    <div className="relative group">
      <div className="w-32 h-32 rounded-full overflow-hidden border-4 border-white shadow-md ring-1">
        {user?.profilePic && (
          <Image
            src={image ? URL.createObjectURL(image) : getFileUrl(user.profilePic)}
            width={200}
            height={200}
            alt="Profile"
            className="w-full h-full object-cover transition-transform group-hover:scale-105"
          />
        )}
        {!user?.profilePic && (
          <div className="w-full h-full flex items-center justify-center bg-gray-200">
            <UserIcon size={32} />
          </div>
        )}
      </div>
      {/* Floating Edit Button */}
      <button
        onClick={handleImageClick}
        className="absolute bottom-1 right-1 p-2 rounded-full border-2 border-white shadow-lg transition-transform hover:scale-110 active:scale-90"
        style={{ backgroundColor: 'oklch(0.51 0.17 280)', color: 'oklch(0.985 0 0)' }}
        disabled={isPending}>
        <Camera size={16} />
      </button>
      <input type="file" ref={fileInputRef} onChange={handleImageChange} className="hidden" accept="image/*" />
    </div>
  )
}
