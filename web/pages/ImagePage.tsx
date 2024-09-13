import {
  Connection,
  Controls,
  MiniMap,
  NodeTypes,
  ReactFlow,
  addEdge,
  useEdgesState,
  useNodesState,
} from '@xyflow/react'
import '@xyflow/react/dist/base.css'
import { useCallback } from 'react'
import { useSearchParams } from 'react-router-dom'

import { Badge } from '../components/Badge'
import CustomNode from '../components/Node'
import { FluentChevronDown20Regular } from '../components/icons/fluent-chevron-down-20-regular'
import { FluentChevronUp20Regular } from '../components/icons/fluent-chevron-up-20-regular'
import { Quay } from '../components/icons/quay'
import { SimpleIconsDocker } from '../components/icons/simple-icons-docker'
import { SimpleIconsGit } from '../components/icons/simple-icons-git'
import { SimpleIconsGithub } from '../components/icons/simple-icons-github'
import { SimpleIconsGitlab } from '../components/icons/simple-icons-gitlab'

interface Tag {
  label: string
  color: string
}

const nodeTypes: NodeTypes = {
  custom: CustomNode,
}

const initNodes = [
  {
    id: '1',
    type: 'custom',
    data: { subtitle: 'default', title: 'Namespace', label: 'N' },
    position: { x: 0, y: 50 },
  },
  {
    id: '2',
    type: 'custom',
    data: { subtitle: 'home-assistant', title: 'Deployment', label: 'D' },
    position: { x: 0, y: 150 },
  },
  {
    id: '3',
    type: 'custom',
    data: { subtitle: 'home-assistant', title: 'Pod', label: 'P' },
    position: { x: 0, y: 250 },
  },
  {
    id: '4',
    type: 'custom',
    data: { subtitle: 'home-assistant', title: 'Container', label: 'C' },
    position: { x: 0, y: 350 },
  },

  {
    id: '8',
    type: 'custom',
    data: { subtitle: 'test', title: 'Namespace', label: 'N' },
    position: { x: 250, y: 150 },
  },
  {
    id: '6',
    type: 'custom',
    data: { subtitle: 'home-assistant', title: 'Pod', label: 'P' },
    position: { x: 250, y: 250 },
  },
  {
    id: '7',
    type: 'custom',
    data: { subtitle: 'home-assistant', title: 'Container', label: 'C' },
    position: { x: 250, y: 350 },
  },

  {
    id: '5',
    type: 'custom',
    data: { subtitle: 'home-assistant', title: 'Image', label: 'I' },
    position: { x: 0, y: 450 },
  },
]

const initEdges = [
  {
    id: 'e1',
    source: '1',
    target: '2',
  },
  {
    id: 'e2',
    source: '2',
    target: '3',
  },
  {
    id: 'e3',
    source: '3',
    target: '4',
  },
  {
    id: 'e4',
    source: '4',
    target: '5',
  },
  {
    id: 'e5',
    source: '6',
    target: '7',
  },
  {
    id: 'e6',
    source: '7',
    target: '5',
  },
  {
    id: 'e7',
    source: '8',
    target: '6',
  },
]

function MockMarkdown(): JSX.Element {
  return (
    <>
      <ul>
        <li>
          Fix BTHome validate triggers for device with multiple buttons (
          <a href="https://github.com/thecode">@thecode</a> -{' '}
          <a
            href="https://github.com/home-assistant/core/pull/125183"
            data-hovercard-type="pull_request"
            data-hovercard-url="/home-assistant/core/pull/125183/hovercard"
          >
            #125183
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/bthome/"
            rel="nofollow"
          >
            bthome docs
          </a>
          )
        </li>
        <li>
          Improve play media support in LinkPlay (
          <a href="https://github.com/silamon">@silamon</a> -{' '}
          <a href="https://github.com/home-assistant/core/pull/125205">
            #125205
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/linkplay/"
            rel="nofollow"
          >
            linkplay docs
          </a>
          )
        </li>
        <li>
          Revert{' '}
          <a href="https://github.com/home-assistant/core/pull/122676">
            #122676
          </a>{' '}
          Yamaha discovery (<a href="https://github.com/joostlek">@joostlek</a>{' '}
          -{' '}
          <a href="https://github.com/home-assistant/core/pull/125216">
            #125216
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/yamaha/"
            rel="nofollow"
          >
            yamaha docs
          </a>
          )
        </li>
        <li>
          Fix blocking call in yale_smart_alarm (
          <a href="https://github.com/gjohansson-ST">@gjohansson-ST</a> -{' '}
          <a href="https://github.com/home-assistant/core/pull/125255">
            #125255
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/yale_smart_alarm/"
            rel="nofollow"
          >
            yale_smart_alarm docs
          </a>
          )
        </li>
        <li>
          Don't show input panel if default code provided in envisalink (
          <a href="https://github.com/gjohansson-ST">@gjohansson-ST</a> -{' '}
          <a href="https://github.com/home-assistant/core/pull/125256">
            #125256
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/envisalink/"
            rel="nofollow"
          >
            envisalink docs
          </a>
          )
        </li>
        <li>
          Increase AquaCell timeout and handle timeout exception properly (
          <a href="https://github.com/Jordi1990">@Jordi1990</a> -{' '}
          <a href="https://github.com/home-assistant/core/pull/125263">
            #125263
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/aquacell/"
            rel="nofollow"
          >
            aquacell docs
          </a>
          )
        </li>
        <li>
          Bump aiorussound to 3.0.4 (
          <a href="https://github.com/noahhusby">@noahhusby</a> -{' '}
          <a href="https://github.com/home-assistant/core/pull/125285">
            #125285
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/russound_rio/"
            rel="nofollow"
          >
            russound_rio docs
          </a>
          )
        </li>
        <li>
          Add follower to the PlayingMode enum (
          <a href="https://github.com/silamon">@silamon</a> -{' '}
          <a href="https://github.com/home-assistant/core/pull/125294">
            #125294
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/linkplay/"
            rel="nofollow"
          >
            linkplay docs
          </a>
          )
        </li>
        <li>
          Don't allow templating min, max, step in config entry template number
          (<a href="https://github.com/emontnemery">@emontnemery</a> -{' '}
          <a href="https://github.com/home-assistant/core/pull/125342">
            #125342
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/template/"
            rel="nofollow"
          >
            template docs
          </a>
          )
        </li>
        <li>
          Fix for Hue sending effect None at turn_on command while no effect is
          active (<a href="https://github.com/marcelveldt">@marcelveldt</a> -{' '}
          <a href="https://github.com/home-assistant/core/pull/125377">
            #125377
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/hue/"
            rel="nofollow"
          >
            hue docs
          </a>
          )
        </li>
        <li>
          Lyric: fixed missed snake case conversions (
          <a href="https://github.com/dalinicus">@dalinicus</a> -{' '}
          <a href="https://github.com/home-assistant/core/pull/125382">
            #125382
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/lyric/"
            rel="nofollow"
          >
            lyric docs
          </a>
          )
        </li>
        <li>
          Bump pysmlight to 0.0.14 (
          <a href="https://github.com/tl-sl">@tl-sl</a> -{' '}
          <a href="https://github.com/home-assistant/core/pull/125387">
            #125387
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/smlight/"
            rel="nofollow"
          >
            smlight docs
          </a>
          )
        </li>
        <li>
          Bump pypck to 0.7.22 (
          <a href="https://github.com/alengwenus">@alengwenus</a> -{' '}
          <a href="https://github.com/home-assistant/core/pull/125389">
            #125389
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/lcn/"
            rel="nofollow"
          >
            lcn docs
          </a>
          )
        </li>
        <li>
          Increase coordinator update_interval for fyta (
          <a href="https://github.com/dontinelli">@dontinelli</a> -{' '}
          <a
            href="https://github.com/home-assistant/core/pull/125393"
            data-hovercard-type="pull_request"
            data-hovercard-url="/home-assistant/core/pull/125393/hovercard"
          >
            #125393
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/fyta/"
            rel="nofollow"
          >
            fyta docs
          </a>
          )
        </li>
        <li>
          Fix controlling AC temperature in airtouch5 (
          <a href="https://github.com/danzel">@danzel</a> -{' '}
          <a
            href="https://github.com/home-assistant/core/pull/125394"
            data-hovercard-type="pull_request"
            data-hovercard-url="/home-assistant/core/pull/125394/hovercard"
          >
            #125394
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/airtouch5/"
            rel="nofollow"
          >
            airtouch5 docs
          </a>
          )
        </li>
        <li>
          Bump sfrbox-api to 0.0.10 (
          <a href="https://github.com/AlexT59">@AlexT59</a> -{' '}
          <a
            href="https://github.com/home-assistant/core/pull/125405"
            data-hovercard-type="pull_request"
            data-hovercard-url="/home-assistant/core/pull/125405/hovercard"
          >
            #125405
          </a>
          )
        </li>
        <li>
          Improve handling of old firmware versions (
          <a href="https://github.com/tl-sl">@tl-sl</a> -{' '}
          <a
            href="https://github.com/home-assistant/core/pull/125406"
            data-hovercard-type="pull_request"
            data-hovercard-url="/home-assistant/core/pull/125406/hovercard"
          >
            #125406
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/smlight/"
            rel="nofollow"
          >
            smlight docs
          </a>
          )
        </li>
        <li>
          Set min_power similar to max_power to support all inverters from
          apsystems (
          <a href="https://github.com/mawoka-myblock">@mawoka-myblock</a> -{' '}
          <a
            href="https://github.com/home-assistant/core/pull/124247"
            data-hovercard-type="pull_request"
            data-hovercard-url="/home-assistant/core/pull/124247/hovercard"
          >
            #124247
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/apsystems/"
            rel="nofollow"
          >
            apsystems docs
          </a>
          )
        </li>
        <li>
          Update frontend to 20240906.0 (
          <a href="https://github.com/piitaya">@piitaya</a> -{' '}
          <a
            href="https://github.com/home-assistant/core/pull/125409"
            data-hovercard-type="pull_request"
            data-hovercard-url="/home-assistant/core/pull/125409/hovercard"
          >
            #125409
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/frontend/"
            rel="nofollow"
          >
            frontend docs
          </a>
          )
        </li>
        <li>
          Bump pyatv to 0.15.1 (
          <a href="https://github.com/postlund">@postlund</a> -{' '}
          <a
            href="https://github.com/home-assistant/core/pull/125412"
            data-hovercard-type="pull_request"
            data-hovercard-url="/home-assistant/core/pull/125412/hovercard"
          >
            #125412
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/apple_tv/"
            rel="nofollow"
          >
            apple_tv docs
          </a>
          )
        </li>
      </ul>
    </>
  )
}

function MockDescription(): JSX.Element {
  return (
    <>
      <h1 id="home-assistant">Home Assistant</h1>
      <p>
        Open source home automation that puts local control and privacy first.
        Powered by a worldwide community of tinkerers and DIY enthusiasts.
        Perfect to run on a Raspberry Pi or a local server.
      </p>
      <p>
        Check out <a href="https://home-assistant.io">home-assistant.io</a> for{' '}
        <a href="https://home-assistant.io/demo/">a demo</a>,{' '}
        <a href="https://home-assistant.io/getting-started/">
          installation instructions
        </a>
        ,
        <a href="https://home-assistant.io/getting-started/automation-2/">
          tutorials
        </a>{' '}
        and <a href="https://home-assistant.io/docs/">documentation</a>.
      </p>
      <p>
        <img
          src="https://raw.github.com/home-assistant/home-assistant/master/docs/screenshots.png"
          alt="screenshot states"
        />
      </p>
      <h2 id="featured-integrations">Featured integrations</h2>
      <p>
        <img
          src="https://raw.github.com/home-assistant/home-assistant/dev/docs/screenshot-components.png"
          alt="screenshot components"
        />
      </p>
      <p>
        If you run into issues while using Home Assistant, check the{' '}
        <a href="https://home-assistant.io/help/">
          Home Assistant help section
        </a>{' '}
        of our website for further help and information.
      </p>
    </>
  )
}

export function ImagePage(): JSX.Element {
  const [params, _] = useSearchParams()

  const imageName = params.get('name')
  const imageVersion = params.get('version')

  const [nodes, setNodes, onNodesChange] = useNodesState(initNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState(initEdges)

  const tags: Tag[] = [
    { label: 'k8s', color: 'bg-blue-100' },
    { label: 'docker', color: 'bg-red-100' },
    { label: 'Pod', color: 'bg-orange-100' },
    { label: 'Job', color: 'bg-blue-100' },
    { label: 'ghcr', color: 'bg-blue-100' },
    { label: 'up-to-date', color: 'bg-green-100' },
    { label: 'outdated', color: 'bg-red-100' },
  ]

  const onConnect = useCallback(
    (connection: Connection) => setEdges((eds) => addEdge(connection, eds)),
    []
  )

  return (
    <div className="flex flex-col items-center w-full py-[40px] px-[20px]">
      {/* Header */}
      <img
        className="w-16 rounded"
        src="https://www.gravatar.com/avatar/461df105cc6cfcf386ebd5af85b925dc?s=120&r=g&d=404"
      />
      <h1 className="text-2xl font-medium">{imageName}</h1>
      <div className="flex items-center">
        <FluentChevronDown20Regular className="text-red-500" />
        <p className="font-medium text-red-500">{imageVersion}</p>
        <p className="font-medium ml-4 text-green-500">{imageVersion}</p>
        <FluentChevronUp20Regular className="text-green-500" />
      </div>
      <div className="flex mt-2 items-center">
        {tags.map((x) => (
          <Badge label={x.label} color={x.color} />
        ))}
      </div>

      {/* Release notes */}
      <div className="flex mt-2 space-x-4 items-center">
        <a
          title="Project's page on GitHub"
          href="https://github.com/home-assistant/core"
          target="_blank"
        >
          <SimpleIconsGithub className="text-black" />
        </a>
        <a
          title="Project's page on GitLab"
          href="https://gitlab.com/arm-research/smarter/smarter-device-manager"
          target="_blank"
        >
          <SimpleIconsGitlab className="text-orange-500" />
        </a>
        <a
          title="Project's page on Docker Hub"
          href="https://hub.docker.com/r/homeassistant/home-assistant"
          target="_blank"
        >
          <SimpleIconsDocker className="text-blue-500" />
        </a>
        <a
          title="Project's page on Quay"
          href="https://quay.io/repository/jetstack/cert-manager-webhook?tab=info"
          target="_blank"
        >
          <Quay className="text-blue-500" />
        </a>
        <a
          title="Project's source code"
          href="https://quay.io/repository/jetstack/cert-manager-webhook?tab=info"
          target="_blank"
        >
          <SimpleIconsGit className="text-orange-500" />
        </a>
      </div>

      <main className="min-w-[200px] max-w-[980px] box-border space-y-6 mt-6">
        <div className="rounded-lg bg-white px-4 py-6 shadow">
          <div className="markdown-body">
            <MockDescription />
          </div>
        </div>

        <div className="rounded-lg bg-white px-4 py-6 shadow">
          <div className="markdown-body">
            <h1>{imageVersion}</h1>
            <MockMarkdown />
          </div>
        </div>

        {/* Graph */}
        <div className="rounded-lg bg-white px-4 py-2 shadow h-[480px]">
          <ReactFlow
            nodes={nodes}
            edges={edges}
            onNodesChange={onNodesChange}
            nodeTypes={nodeTypes}
            fitView
            edgesFocusable={false}
            nodesDraggable={true}
            nodesConnectable={false}
            nodesFocusable={false}
            draggable={true}
            panOnDrag={true}
            elementsSelectable={false}
          >
            <Controls />
          </ReactFlow>
        </div>
      </main>
    </div>
  )
}
