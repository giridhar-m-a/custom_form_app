import { FieldType } from '@/types/form.types'
import { IconType } from 'react-icons/lib'
import {
  MdAccessTime,
  MdArrowDropDownCircle,
  MdAudioFile,
  MdCheckBox,
  MdDateRange,
  MdEmail,
  MdFolder,
  MdImage,
  MdLink,
  MdList,
  MdNotes,
  MdNumbers,
  MdOutlineColorLens,
  MdPermContactCalendar,
  MdPhone,
  MdRadioButtonChecked,
  MdShortText,
  MdStar,
  MdVideoFile
} from 'react-icons/md'
import { RxSlider } from 'react-icons/rx'

export const FIELD_TYPE_OPTIONS: { value: FieldType; label: string; icon?: IconType }[] = [
  { value: 'text', label: 'Text', icon: MdShortText },
  { value: 'textArea', label: 'Textarea', icon: MdNotes },
  { value: 'number', label: 'Number', icon: MdNumbers },
  { value: 'email', label: 'Email', icon: MdEmail },
  { value: 'phone', label: 'Phone', icon: MdPhone },
  { value: 'url', label: 'URL', icon: MdLink },
  { value: 'date', label: 'Date', icon: MdDateRange },
  { value: 'time', label: 'Time', icon: MdAccessTime },
  { value: 'datetime', label: 'Date & Time', icon: MdPermContactCalendar },

  { value: 'file', label: 'File Upload', icon: MdFolder },
  { value: 'image', label: 'Image Upload', icon: MdImage },
  { value: 'video', label: 'Video Upload', icon: MdVideoFile },
  { value: 'audio', label: 'Audio Upload', icon: MdAudioFile },

  { value: 'checkbox', label: 'Checkbox', icon: MdCheckBox },
  { value: 'radio', label: 'Radio', icon: MdRadioButtonChecked },
  { value: 'dropdown', label: 'Single Select', icon: MdArrowDropDownCircle },
  { value: 'multiselect', label: 'Multi Select', icon: MdList },

  { value: 'rating', label: 'Rating', icon: MdStar },
  { value: 'slider', label: 'Slider', icon: RxSlider },
  { value: 'color', label: 'Color Picker', icon: MdOutlineColorLens }
]
